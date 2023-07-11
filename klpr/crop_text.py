import os
import sys
import unicodedata
import numpy as np

import cv2 as cv

from ch07.vehicle_no_2 import ChainParam, DetectCustom, DetectYolo
import easyocr


def main():
    # root = '/Users/kcson/mywork/data/[원천]자동차번호판OCR데이터'
    root = '../imgs/car'
    train_file_list = os.listdir(root)

    index = 0
    valid_index = 0
    h_sum = 0
    w_sum = 0
    for file_name in train_file_list:
        if index == 50:
            break
        file_name = unicodedata.normalize('NFC', file_name)
        label = os.path.splitext(file_name)[0].split('-')[0]
        # if len(label) > 8:
        #     continue

        full_path = os.path.join(root, file_name)
        print(full_path)
        crop_text(full_path)
        index += 1


def crop_text(full_path):
    src = cv.imread(full_path)
    # src = cv.imread('../imgs/vehicle23.jpeg')
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
        # sys.exit()
        return

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

    kernel = cv.getStructuringElement(cv.MORPH_RECT, (6, 6))
    # src_bin = cv.erode(src_bin, kernel, iterations=1)
    # src_bin = cv.dilate(src_bin, kernel, iterations=1)

    contours, _ = cv.findContours(src_bin, cv.RETR_LIST, cv.CHAIN_APPROX_NONE)
    boxes = []
    box_contours = []
    box_areas = []
    for contour in contours:
        x, y, w, h = cv.boundingRect(contour)
        if not is_valid_contour(contour, contours, src):
            continue
        boxes.append((x, y, w, h))

    boxes = remove_inner_box(boxes, src_bin)
    # boxes = remove_invalid_box(boxes)

    sorted_boxes = []
    boxes = sorted(boxes, key=lambda box: box[0])
    boxes = sort_box(boxes, sorted_boxes)

    boxes = merge_box(boxes)
    for i, box in enumerate(boxes, start=1):
        x, y, w, h = box
        box_areas.append(h * w)
        print('area : ', h * w, ' ratio : ', w / h)
        cv.rectangle(src, (x, y, w, h), (0, 0, 255), thickness=1)
        cv.putText(src, str(i), (x, y), cv.FONT_HERSHEY_SIMPLEX, 1, (0, 0, 255), 2, cv.LINE_AA)

    print('average : ', np.mean(box_areas))
    print('std     : ', np.std(box_areas))
    # recognition_text(boxes, src_bin)
    cv.imshow('src_bin', src_bin)
    cv.imshow('src', src)
    cv.waitKey()
    cv.destroyAllWindows()


def merge_box(boxes):
    return boxes


def merge(a, b):
    x = min(a[0], b[0])
    y = min(a[1], b[1])
    w = max(a[0] + a[2], b[0] + b[2]) - x
    h = max(a[1] + a[3], b[1] + b[3]) - y
    return [x, y, w, h]


def sort_box(boxes, sorted_boxes):
    for box1 in boxes:
        if box1 in sorted_boxes:
            continue
        is_next_box = True
        for box2 in boxes:
            if box2 in sorted_boxes:
                continue
            if not get_is_next_box(box1, box2, boxes):
                is_next_box = False
                break
        if not is_next_box:
            continue
        sorted_boxes.append(box1)
        break

    if len(boxes) == len(sorted_boxes):
        return sorted_boxes

    return sort_box(boxes, sorted_boxes)


def get_is_next_box(box1, box2, boxes) -> bool:
    x1, y1, w1, h1 = box1
    x2, y2, w2, h2 = box2
    if y1 > y2 + h2:
        for box in boxes:
            if is_overlap_y(box1, box) and is_overlap_y(box, box2):
                if not is_overlap_x(box1, box2):
                    return True
                else:
                    return False
        return False
    if y1 + h1 < y2:
        return True
    if x1 + w1 < x2:
        return True
    if x1 + w1 < x2 + w2:
        return True
    if x1 + w1 >= x2 + w2 and y1 > y2:
        return False
    if x1 + w1 >= x2 + w2 and y1 <= y2:
        return True

    return True


def is_overlap_y(box1, box2) -> bool:
    x1, y1, w1, h1 = box1
    x2, y2, w2, h2 = box2
    if y2 <= y1 <= y2 + h2 or y2 <= y1 + h1 <= y2 + h2:
        return True

    if y1 <= y2 <= y1 + h1 or y1 <= y2 + h2 <= y1 + h1:
        return True

    return False


def is_overlap_x(box1, box2) -> bool:
    x1, y1, w1, h1 = box1
    x2, y2, w2, h2 = box2

    if x2 <= x1 <= x2 + w2 or x2 <= x1 + w1 <= x2 + w2:
        return True

    if x1 <= x2 <= x1 + w1 or x1 <= x2 + w2 <= x1 + w1:
        return True
    return False


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
    # approx = cv.approxPolyDP(contour, 10, True)
    # cv.drawContours(src, approx, 0, (0, 0, 255), thickness=3)
    # cv.rectangle(src, cv.boundingRect(approx), (0, 0, 255), thickness=1)
    # if not cv.isContourConvex(approx):
    #     return False

    new_height, new_width = src.shape[:2]
    x, y, w, h = cv.boundingRect(contour)
    if w / h > 6 or w / h < 0.15:
        return False
    if w * h < 1200:
        return False
    if w / new_width > 0.9 or h / new_height > 0.9:
        return False
    return True


def remove_invalid_box(boxes):
    boxes = sorted(boxes, key=lambda box: box[2] * box[3], reverse=True)
    max_area = boxes[0][2] * boxes[0][3]

    valid_boxes = []
    for box1 in boxes:
        x1, y1, w1, h1 = box1
        if (w1 * h1) / max_area < 0.1:
            continue
        valid_boxes.append(box1)

    return valid_boxes


def remove_inner_box(boxes, src_bin):
    outer_boxes = []
    temp_boxes = []
    for box1 in boxes:
        x1_min, y1_min, w1, h1 = box1
        target = src_bin[y1_min:y1_min + h1, x1_min:x1_min + w1]
        contours, _ = cv.findContours(target, cv.RETR_LIST, cv.CHAIN_APPROX_NONE)
        if len(contours) > 5:
            continue
        temp_boxes.append(box1)

    for box1 in temp_boxes:
        is_inner_box = False
        x1_min, y1_min, w1, h1 = box1
        x1_max = x1_min + w1
        y1_max = y1_min + h1
        for box2 in temp_boxes:
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
    # crop_text('')
