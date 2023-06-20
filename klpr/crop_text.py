import sys

import cv2 as cv

from ch07.vehicle_no_2 import ChainParam, DetectCustom, DetectYolo
import easyocr


def main():
    src = cv.imread('../imgs/vehicle20.jpeg')
    if src is None:
        print('image load fail!!')
        sys.exit()

    chainParam = ChainParam(src)
    chainParam.copySrc = src.copy()

    detect = DetectYolo()
    detect.set_next(DetectCustom())
    detect.handle(chainParam)

    if chainParam.dst is None:
        print('chainParam.dst is Non')
        sys.exit()

    src = chainParam.dst
    height, width = src.shape[:2]
    ratio = width / height
    new_width = 1024
    new_height = int(new_width / ratio)
    src = cv.resize(chainParam.dst, (new_width, new_height))

    src_bin = cv.cvtColor(src, cv.COLOR_BGR2GRAY)
    src_bin = cv.GaussianBlur(src_bin, (0, 0), 2, sigmaY=0, borderType=cv.BORDER_DEFAULT)
    # src_bin = cv.adaptiveThreshold(src_bin, 255, cv.ADAPTIVE_THRESH_GAUSSIAN_C, cv.THRESH_BINARY_INV, 19, 4)
    thres, _ = cv.threshold(src_bin, 0, 255, cv.THRESH_BINARY_INV | cv.THRESH_OTSU, dst=src_bin)

    kernel = cv.getStructuringElement(cv.MORPH_RECT, (5, 5))
    # src_bin = cv.erode(src_bin, kernel)
    # src_bin = cv.dilate(src_bin, kernel)

    contours, _ = cv.findContours(src_bin, cv.RETR_LIST, cv.CHAIN_APPROX_NONE)
    boxes = []
    for contour in contours:
        x, y, w, h = cv.boundingRect(contour)
        if not is_valid_contour(contour, contours, src):
            continue
        boxes.append((x, y, w, h))

    boxes = remove_inner_box(boxes)
    boxes = sort_box(boxes)
    boxes = merge_box(boxes)
    for i, box in enumerate(boxes,start=1):
        x, y, w, h = box
        cv.rectangle(src, (x, y, w, h), (0, 0, 255), thickness=1)
        cv.putText(src, str(i), (x, y), cv.FONT_HERSHEY_SIMPLEX, 1, (0, 0, 255), 2, cv.LINE_AA)

    # recognition_text(boxes, src_bin)
    cv.imshow('src_bin', src_bin)
    cv.imshow('src', src)
    cv.waitKey()
    cv.destroyAllWindows()


def merge_box(boxes):
    return boxes


def sort_box(boxes):
    sorted_box = []
    boxes = sorted(boxes, key=lambda box: box[0])
    return boxes


def recognition_text(boxes, src_bin):
    white_list = '1234567890가나다라마거너더러머버서어저고노도로모보소오조구누두루무부수우주아바사자하허호배국합육해공인천대'
    reader = easyocr.Reader(['ko'])
    for box in boxes:
        x, y, w, h = box
        target = src_bin[y - 10:y + h + 10, x - 10:x + w + 10]
        cv.imshow('target', target)
        cv.waitKey()
        result = reader.readtext(target, detail=0, paragraph=True, allowlist=white_list)
        for v in result:
            print(v)


def is_valid_contour(contour, contours, src) -> bool:
    new_height, new_width = src.shape[:2]
    x, y, w, h = cv.boundingRect(contour)
    if w / h > 6 or w / h < 0.1:
        return False
    if w * h < 1400:
        return False
    if w / new_width > 0.9 or h / new_height > 0.9:
        return False
    return True


def remove_inner_box(boxes):
    outer_boxes = []
    for box1 in boxes:
        is_inner_box = False
        x1_min, y1_min, w1, h1 = box1
        x1_max = x1_min + w1
        y1_max = y1_min + h1
        for box2 in boxes:
            x2_min, y2_min, w2, h2 = box2
            x2_max = x2_min + w2
            y2_max = y2_min + h2
            if x1_min > x2_min and x1_max < x2_max and y1_min > y2_min and y1_max < y2_max:
                is_inner_box = True
                break
        if not is_inner_box:
            outer_boxes.append(box1)
    return outer_boxes


if __name__ == "__main__":
    main()
