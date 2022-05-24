from __future__ import annotations

import random
from abc import ABC, abstractmethod
from dataclasses import dataclass
import cv2 as cv
import sys

import numpy as np


@dataclass
class ChainParam:
    src: np
    pre: np = None
    dst: np = None
    vehicle_no: str = ''


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
        param.pre = src_bin

        return super().handle(param)


class CutVehicleRegion(AbstractHandler):
    src = None
    pre = None

    def handle(self, param: ChainParam):
        self.src = param.src
        self.pre = param.pre
        min_pts = None
        min_area: int = sys.maxsize
        contours, _ = cv.findContours(param.pre, cv.RETR_LIST, cv.CHAIN_APPROX_NONE)
        for pts in contours:
            x, y, w, h = cv.boundingRect(pts)
            if h / w > 0.5:
                continue
            area = w * h
            if area < 300 * 250:
                continue

            count = self.get_connected_components((x, y, w, h))
            if count < 4 or count > 20:
                continue

            if area < min_area:
                print(area)
                min_area = area
                min_pts = pts

        if min_pts is not None:
            x, y, w, h = cv.boundingRect(min_pts)
            param.dst = param.src[y:y + h, x:x + w]
            cv.rectangle(param.src, (x, y), (x + w, y + h), (0, 0, 255), thickness=3)

        return super().handle(param)

    def get_connected_components(self, rect) -> int:
        component_count = 0
        x, y, w, h = rect
        mat = self.pre[y:y + h, x:x + w]
        contours, _ = cv.findContours(mat, cv.RETR_LIST, cv.CHAIN_APPROX_NONE)
        if len(contours) == 0:
            return 0
        for pts in contours:
            rect = cv.boundingRect(pts)
            if cv.contourArea(pts) < 200:
                continue
            if rect[2] * rect[3] < 200:
                continue
            if not self.guess_vehicle_no(pts):
                continue

            similar_count = self.similar_contour(contours, pts)
            if similar_count < 3:
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
        for pts in contours:
            x2, y2, w2, h2 = cv.boundingRect(pts)
            area_diff: float = abs(w1 * h1 - w2 * h2) / (w1 * h1)
            width_diff: float = abs(w1 - w2) / w1
            height_diff: float = abs(h1 - h2) / h1

            if area_diff > 0.6 or width_diff > 0.8 or height_diff > 0.2:
                continue

            similar_count = similar_count + 1
        return similar_count


class GetVehicleNo(AbstractHandler):
    def handle(self, param: ChainParam):
        print('GetVehicleNo')
        return super().handle(param)


if __name__ == "__main__":
    print("main")
    src = cv.imread('../imgs/vehicle.jpeg')
    if src is None:
        print('image read fail!!')
        sys.exit()

    chainParam = ChainParam(src)
    # chainParam.src = src

    pre = PreProcess()
    pre.set_next(CutVehicleRegion()).set_next(GetVehicleNo())
    pre.handle(chainParam)

    cv.imshow('src', src)
    cv.imshow('pre', chainParam.pre)
    if chainParam.dst is not None:
        cv.imshow('dst', chainParam.dst)

    cv.waitKey()
    cv.destroyAllWindows()
