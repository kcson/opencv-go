from __future__ import annotations

import math
import random
import re
from abc import ABC, abstractmethod
from dataclasses import dataclass
import cv2 as cv
import sys

import easyocr
import numpy as np
import pytesseract
from PIL import ImageFont, ImageDraw, Image
import imutils


@dataclass
class ChainParam:
    src: np
    copySrc: np = None
    pre: np = None
    dst: np = None
    vehicle_no: str = ''
    detectedLP = None


class Handler(ABC):
    @abstractmethod
    def set_next(self, handler: Handler) -> Handler:
        pass

    @abstractmethod
    def handle(self, param: ChainParam):
        pass


class AbstractHandler(Handler):
    _next_handler: Handler = None

    def set_next(self, handler: Handler) -> Handler:
        self._next_handler = handler

        return handler

    def handle(self, param: ChainParam):
        if self._next_handler:
            return self._next_handler.handle(param)

        return None


class PreProcess(AbstractHandler):
    def handle(self, param: ChainParam):
        gray = cv.cvtColor(param.src, cv.COLOR_BGR2GRAY)
        src_bin = cv.GaussianBlur(gray, (0, 0), 2, sigmaY=0, borderType=cv.BORDER_DEFAULT)
        src_bin = cv.adaptiveThreshold(src_bin, 255, cv.ADAPTIVE_THRESH_GAUSSIAN_C, cv.THRESH_BINARY_INV, 19, 4)

        # src_bin = cv.erode(src_bin, None, iterations=2)
        src_bin = cv.dilate(src_bin, None, iterations=1)

        param.pre = src_bin

        return super().handle(param)


class DetectYolo(AbstractHandler):
    def handle(self, param: ChainParam):
        if param.dst is not None:
            return super().handle(param)
        param.src = param.copySrc.copy()

        tm = cv.TickMeter()
        tm.start()
        net = cv.dnn.readNetFromDarknet('../klpr/yolov4-ANPR.cfg', '../klpr/yolov4-ANPR.weights')
        tm.stop()
        print('DetectYolo readNetFromDarknet : ', tm.getTimeSec())
        with open('../klpr/obj.names', 'r') as f:
            classes = [line.strip() for line in f.readlines()]

        layer_names = net.getLayerNames()
        output_layers = [layer_names[i - 1] for i in net.getUnconnectedOutLayers()]

        height, width, channels = param.src.shape
        param.src = cv.resize(param.src, None, fx=416 / width, fy=416 / height)
        blob = cv.dnn.blobFromImage(param.src, 0.00392, (416, 416), (0, 0, 0), True, crop=False)
        net.setInput(blob)

        tm = cv.TickMeter()
        tm.start()
        outs = net.forward(output_layers)
        tm.stop()
        print('DetectYolo forward : ', tm.getTimeSec())

        class_ids = []
        confidences = []
        boxes = []

        tm = cv.TickMeter()
        tm.start()
        for out in outs:
            for detection in out:
                scores = detection[5:]
                class_id = np.argmax(scores)
                confidence = scores[class_id]

                if confidence > 0.75:
                    center_x = int(detection[0] * width)
                    center_y = int(detection[1] * height)
                    w = int(detection[2] * width)
                    h = int(detection[3] * height)

                    # 객체의 사각형 테두리 중 좌상단 좌표값 찾기
                    x = int(center_x - w / 2)
                    y = int(center_y - h / 2)

                    boxes.append([x, y, w, h])
                    confidences.append(float(confidence))
                    class_ids.append(class_id)

        indexes = cv.dnn.NMSBoxes(boxes, confidences, 0.75, 0.4)

        max_box_area = 0
        max_box = None
        # tm.start()
        for i in range(len(boxes)):
            if i in indexes:
                class_name = classes[class_ids[i]]
                if class_name == 'car':  # and len(indexes) > 1:
                    continue
                x, y, w, h = boxes[i]
                if w * h > max_box_area:
                    max_box = boxes[i]

        if max_box is not None:
            mx, my, mw, mh = max_box
            mx = 10 if mx < 0 else mx
            my = 10 if my < 0 else my
            param.detectedLP = (mx, my, mw, mh)
            param.dst = param.copySrc[my:my + mh, mx:mx + mw]
        tm.stop()
        print('DetectYolo detection : ', tm.getTimeSec())

        return super().handle(param)


class DetectMy(AbstractHandler):
    def handle(self, param: ChainParam):
        if param.dst is not None:
            return super().handle(param)
        param.src = param.copySrc.copy()

        gray = cv.cvtColor(param.src, cv.COLOR_BGR2GRAY)
        bfilter = cv.bilateralFilter(gray, 11, 11, 17)
        edged = cv.Canny(bfilter, 0, 80)
        edged = cv.dilate(edged, None, iterations=1)

        cv.imshow('edged', edged)

        contours, _ = cv.findContours(edged.copy(), cv.RETR_LIST, cv.CHAIN_APPROX_SIMPLE)
        # contours = imutils.grab_contours(contours)
        contours = sorted(contours, key=cv.contourArea, reverse=True)

        location = None
        for contour in contours:
            cv.rectangle(param.copySrc, cv.boundingRect(contour), (0, 0, 255), thickness=3)
            approx = cv.approxPolyDP(contour, 10, True)
            if len(approx) == 4:
                location = cv.boundingRect(contour)
                break

        if location is not None:
            (x, y, w, h) = location
            param.detectedLP = location
            param.dst = param.copySrc[y:y + h, x:x + w]
            # cv.rectangle(param.copySrc, (x, y, w, h), (0, 0, 255), thickness=1)

        return super().handle(param)


class DetectCustom(AbstractHandler):
    def handle(self, param: ChainParam):
        if param.dst is not None:
            return super().handle(param)
        param.src = param.copySrc.copy()

        gray = cv.cvtColor(param.src, cv.COLOR_BGR2GRAY)
        src_bin = cv.GaussianBlur(gray, (0, 0), 2, sigmaY=0, borderType=cv.BORDER_DEFAULT)
        src_bin = cv.adaptiveThreshold(src_bin, 255, cv.ADAPTIVE_THRESH_GAUSSIAN_C, cv.THRESH_BINARY_INV, 19, 4)
        param.pre = cv.dilate(src_bin, None, iterations=1)

        min_pts = None
        min_area: int = sys.maxsize
        contours, _ = cv.findContours(param.pre, cv.RETR_LIST, cv.CHAIN_APPROX_NONE)
        for pts in contours:
            # 외곽선 근사화
            # approx = cv.approxPolyDP(pts, cv.arcLength(pts, True) * 0.02, True)

            # 컨벡스가 아니고, 사각형이 아니면 무시
            # if not cv.isContourConvex(approx):
            #     continue
            # if len(approx) != 4:
            #     continue

            x, y, w, h = cv.boundingRect(pts)
            # cv.rectangle(param.copySrc, (x, y), (x + w, y + h), (0, 0, 255), thickness=1)
            if h / w > 0.5:
                continue
            area = cv.contourArea(pts)
            print(area)
            if area < 80 * 80:
                continue
            # cv.drawContours(param.src, [approx], 0, (0, 0, 255), 3)
            # cv.rectangle(param.copySrc, (x, y), (x + w, y + h), (0, 0, 255), thickness=3)
            # cv.polylines(param.src, [approx], True, (0, 0, 255), 2, cv.LINE_AA)
            count = self.get_connected_components((x, y, w, h), param)
            print('count : ', count)
            if count < 5 or count > 15:
                continue
            #  cv.rectangle(param.copySrc, (x, y), (x + w, y + h), (0, 0, 255), thickness=1)

            if w * h < min_area:
                min_area = w * h
                min_pts = pts

        if min_pts is not None:
            x, y, w, h = cv.boundingRect(min_pts)
            param.detectedLP = (x, y, w, h)
            param.dst = param.copySrc[y:y + h, x + 20:x + w - 20]

        return super().handle(param)

    def get_connected_components(self, rect, param: ChainParam) -> int:
        component_count = 0
        x, y, w, h = rect
        mat = param.pre[y:y + h, x:x + w]
        contours, _ = cv.findContours(mat, cv.RETR_LIST, cv.CHAIN_APPROX_NONE)
        if len(contours) == 0:
            return 0
        for pts in contours:
            if cv.contourArea(pts) < 500:
                continue
            if not self.guess_vehicle_no(pts):
                continue

            similar_count = self.similar_contour(contours, pts)
            if similar_count < 5:
                continue

            component_count = component_count + 1

        return component_count

    @staticmethod
    def guess_vehicle_no(pts) -> bool:
        x, y, w, h = cv.boundingRect(pts)
        ratio = h / w
        if ratio < 1.0 or ratio > 4.0:
            return False
        return True

    @staticmethod
    def similar_contour(contours, pts) -> int:
        similar_count = 0
        x1, y1, w1, h1 = cv.boundingRect(pts)
        diagonal_length = math.sqrt(w1 * w1 + h1 * h1)
        for pts in contours:
            x2, y2, w2, h2 = cv.boundingRect(pts)
            area_diff: float = abs(w1 * h1 - w2 * h2) / (w1 * h1)
            width_diff: float = abs(w1 - w2) / w1
            height_diff: float = abs(h1 - h2) / h1

            if area_diff > 0.6 or width_diff > 0.8 or height_diff > 0.2:
                continue

            distance = math.sqrt((x1 - x2) ** 2 + (y1 - y2) ** 2)
            if distance > diagonal_length * 5:
                continue

            angle_diff = 0
            dx = abs(x1 - x2)
            dy = abs(y1 - y2)
            if dx == 0:
                angle_diff = 90
            else:
                angle_diff = math.atan(dy / dx) * (180 / math.pi)
            if angle_diff > 20.0:
                continue

            similar_count = similar_count + 1
        return similar_count


class RecognitionEasy(AbstractHandler):
    def handle(self, param: ChainParam):
        if param.dst is None:
            return super().handle(param)

        tm = cv.TickMeter()
        tm.start()
        reader = easyocr.Reader(['ko'])
        tm.stop()
        print('RecognitionEasy Reader', tm.getTimeSec())

        tm = cv.TickMeter()
        tm.start()
        height, width = param.dst.shape[:2]
        ratio = width / height
        new_width = 1024
        new_height = int(new_width / ratio)
        param.dst = cv.resize(param.dst, (new_width, new_height))

        src_bin = cv.cvtColor(param.dst, cv.COLOR_BGR2GRAY)
        src_bin = cv.GaussianBlur(src_bin, (0, 0), 1, sigmaY=0, borderType=cv.BORDER_DEFAULT)
        # src_bin = cv.bilateralFilter(src_bin, 0, 10, 5)
        thres, _ = cv.threshold(src_bin, 0, 255, cv.THRESH_BINARY_INV | cv.THRESH_OTSU, dst=src_bin)

        # src_bin = cv.dilate(src_bin, None, iterations=1)
        # src_bin = cv.erode(src_bin, None, iterations=1)

        # src_bin = cv.copyMakeBorder(src_bin, 10, 10, 10, 10, cv.BORDER_CONSTANT, value=(0, 0, 0))
        cv.imshow('src_bin', src_bin)
        result = reader.readtext(src_bin, detail=0, paragraph=True)
        vehicle_no = ''
        for v in result:
            vehicle_no = vehicle_no + v

        result_vehicle_no = ''
        for i in range(len(vehicle_no)):
            v = vehicle_no[i]
            print(v)
            if ('가' <= v <= '힣') or v.isdigit():
                result_vehicle_no = result_vehicle_no + v

        p = re.compile('[0-9]{2,3}[가-힣]{1}[0-9]{4}|[가-힣]{2}[0-9]{2}[가-힣]{1}[0-9]{4}')
        m = p.match(result_vehicle_no)
        if m is not None:
            param.vehicle_no = m.group()
        # param.vehicle_no = result_vehicle_no
        tm.stop()
        print('RecognitionEasy Recognition', tm.getTimeSec())
        return super().handle(param)


class RecognitionTesseract(AbstractHandler):
    def handle(self, param: ChainParam):
        if param.dst is None:
            return super().handle(param)

        p = re.compile('[0-9]{2,3}[가-힣]{1}[0-9]{4}|[가-힣]{2}[0-9]{2}[가-힣]{1}[0-9]{4}')
        gray = cv.cvtColor(param.dst, cv.COLOR_BGR2GRAY)
        src_bin = cv.GaussianBlur(gray, (0, 0), 1, sigmaY=0, borderType=cv.BORDER_DEFAULT)
        thres, _ = cv.threshold(src_bin, 0, 255, cv.THRESH_BINARY_INV | cv.THRESH_OTSU, dst=src_bin)

        src_bin = cv.erode(src_bin, None, iterations=1)

        # src_bin = cv.copyMakeBorder(src_bin, 10, 10, 10, 10, cv.BORDER_CONSTANT, value=(0, 0, 0))
        vehicle_no = self.get_text_from_image(src_bin)
        if vehicle_no != '':
            m = p.match(vehicle_no)
            if m is not None:
                param.vehicle_no = m.group()
                return super().handle(param)

        cv.imshow('src_bin', src_bin)
        # cv.waitKey()
        vehicle_no = self.retry_vehicle_no(gray, thres, p)
        param.vehicle_no = vehicle_no
        print(vehicle_no)
        if vehicle_no == '' or len(vehicle_no) == 4:
            v = self.retry_vehicle_no(gray, thres, p, make_border=True)
            print(v)
            if vehicle_no == '' or len(v) >= 4:
                param.vehicle_no = v

        return super().handle(param)

    def retry_vehicle_no(self, gray, thres, p, iterations=10, make_border=False) -> str:
        for i in range(iterations):
            thres = thres + 10
            src_bin = cv.GaussianBlur(gray, (0, 0), 1, sigmaY=4, borderType=cv.BORDER_DEFAULT)
            thres, _ = cv.threshold(src_bin, thres, 255, cv.THRESH_BINARY, dst=src_bin)

            src_bin = cv.erode(src_bin, None, iterations=2)

            if make_border:
                src_bin = cv.copyMakeBorder(src_bin, 10, 10, 10, 10, cv.BORDER_CONSTANT, value=(0, 0, 0))
            vehicle_no = self.get_text_from_image(src_bin)
            if vehicle_no != '':
                m = p.match(vehicle_no)
                if m is not None:
                    return m.group()
                elif len(vehicle_no) > 4:
                    v = vehicle_no[-4:]
                    if v.isdigit():
                        return v

        return ''

    @staticmethod
    def get_text_from_image(src_bin) -> str:
        result: str = ''
        vehicle_no = pytesseract.image_to_string(src_bin, lang='kor', config='--oem 3 --psm 11')
        for i in range(len(vehicle_no)):
            v = vehicle_no[i]
            print(v)
            if ('가' <= v <= '힣') or v.isdigit():
                result = result + v

        return result


if __name__ == "__main__":
    # url = 'https://parkingcone.s3.ap-northeast-2.amazonaws.com/real/user_vehicle/2023/04/eaee1dfc75754d44046cfb9177c59c8b/1682679738_FVBUUC/f9c2755a3aa4cc06404b8e64f030c'
    # url = 'https://parkingcone.s3.ap-northeast-2.amazonaws.com/real/user_vehicle/2023/05/eaee1dfc75754d44046cfb9177c59c8b/1682679738_FVBUUC/7cccd183e821468e65b81b734'
    # url = 'https://parkingcone.s3.ap-northeast-2.amazonaws.com/real/user_vehicle/2023/04/eaee1dfc75754d44046cfb9177c59c8b/1682679738_FVBUUC/2cbdb867245618e515a61776b26dcd'
    # url = 'https://parkingcone.s3.ap-northeast-2.amazonaws.com/real/user_vehicle/2023/05/bc8e73ed52dc49fe9bf95149b00a9f31/1683169300_DGSPYV/66d64a516a3950a1686ce524453b9d81'
    # image_array = np.asarray(bytearray(requests.get(url).content), dtype=np.uint8)
    # src = cv.imdecode(image_array, cv.IMREAD_COLOR)
    src = cv.imread('../imgs/vehicle18.jpeg')
    if src is None:
        print('image read fail!!')
        sys.exit()

    chainParam = ChainParam(src)
    chainParam.copySrc = src.copy()

    pre = DetectYolo()
    pre.set_next(DetectCustom()).set_next(RecognitionEasy())
    # pre.set_next(RecognitionTesseract())
    pre.handle(chainParam)

    # cv.imshow('pre', chainParam.pre)
    if chainParam.dst is not None:
        cv.imshow('dst', chainParam.dst)
    #     cv.imshow('resize', cv.resize(chainParam.dst, (94, 24)))

    if chainParam.dst is not None:
        cv.rectangle(chainParam.copySrc, chainParam.detectedLP, (0, 0, 255), thickness=3)
    if chainParam.vehicle_no != '':
        x, y, w, h = chainParam.detectedLP

        chainParam.copySrc = Image.fromarray(chainParam.copySrc)
        draw = ImageDraw.Draw(chainParam.copySrc)
        font = ImageFont.truetype("AppleGothic.ttf", 50)
        draw.text((x, y - 50), chainParam.vehicle_no, font=font, fill=(0, 255, 255))
        chainParam.copySrc = np.array(chainParam.copySrc)
        print(chainParam.vehicle_no)

    cv.imshow('src', chainParam.copySrc)
    cv.waitKey()
    cv.destroyAllWindows()
