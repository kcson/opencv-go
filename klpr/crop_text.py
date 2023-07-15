import os
import sys
import unicodedata
import numpy as np

import cv2 as cv

from ch07.vehicle_no_2 import ChainParam, DetectCustom, DetectYolo
import easyocr


def main():
    # root = '/Users/kcson/mywork/data/[원천]자동차번호판OCR데이터'
    # root = '../imgs/car'
    root = '/Users/kcson/mywork/data/lpr_pre'
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
        if os.path.getsize(full_path) < 30 * 1024:
            continue
        print(full_path)
        if crop_text(full_path):
            index += 1

    print('index : ', index)


def crop_text(full_path):
    src = cv.imread(full_path)
    # src = cv.imread('../imgs/car/vehicle25.jpeg')
    if src is None:
        print('image load fail!!')
        sys.exit()

    # height, width = src.shape[:2]
    # ratio = width / height
    # new_width = 1024
    # new_height = int(new_width / ratio)
    # src = cv.resize(src, (new_width, new_height))

    chainParam = ChainParam(src)
    chainParam.copySrc = src.copy()

    # detect = DetectYolo()
    # detect.set_next(DetectCustom())
    detect = DetectCustom()
    detect.set_next(DetectYolo())
    detect.handle(chainParam)

    # chainParam.dst = src
    if chainParam.dst is None:
        print('chainParam.dst is Non')
        # sys.exit()
        return False

    src = chainParam.dst
    height, width = src.shape[:2]
    ratio = width / height
    new_width = 1024
    new_height = int(new_width / ratio)
    src = cv.resize(chainParam.dst, (new_width, new_height))

    src_bin = cv.cvtColor(src, cv.COLOR_BGR2GRAY)
    src_bin = cv.GaussianBlur(src_bin, (0, 0), 3, sigmaY=0, borderType=cv.BORDER_DEFAULT)
    # src_bin = cv.adaptiveThreshold(src_bin, 255, cv.ADAPTIVE_THRESH_GAUSSIAN_C, cv.THRESH_BINARY_INV, 19, 4)
    thres, _ = cv.threshold(src_bin, 0, 255, cv.THRESH_BINARY_INV | cv.THRESH_OTSU, dst=src_bin)

    kernel = cv.getStructuringElement(cv.MORPH_RECT, (6, 6))
    # src_bin = cv.dilate(src_bin, kernel, iterations=1)
    # src_bin = cv.erode(src_bin, kernel, iterations=2)

    contours, _ = cv.findContours(src_bin, cv.RETR_LIST, cv.CHAIN_APPROX_NONE)

    # contour 면적이 너무 작거나 큰 contour 제거
    contours = check_area(contours, area_thres_min=750, area_thres_max=60000)
    # 내부에 있는 contour 제거
    contours = remove_inner_box(contours)
    # contour 정렬(위->아래, 좌->우)
    contours = sort_contour(contours)

    # boxes = []
    # box_contours = []
    # for i, contour in enumerate(contours):
    #     x, y, w, h = cv.boundingRect(contour)
    #     if not is_valid_contour(i, contour, contours, src):
    #         continue
    #     boxes.append((x, y, w, h))
    #
    # boxes = remove_inner_box(boxes, src_bin)
    # boxes = sort_box(boxes)
    # boxes = remove_invalid_box(boxes)
    # boxes = merge_box(boxes)
    for i, contour in enumerate(contours, start=1):
        x, y, w, h = cv.boundingRect(contour)
        print('area : ', h * w, ' ratio : ', w / h)
        cv.rectangle(src, (x, y, w, h), (0, 0, 255), thickness=1)
        cv.putText(src, str(i), (x, y), cv.FONT_HERSHEY_SIMPLEX, 1, (0, 0, 255), 2, cv.LINE_AA)

    # recognition_text(boxes, src_bin)
    cv.imshow('src_bin', src_bin)
    cv.imshow('src', src)
    cv.waitKey()
    cv.destroyAllWindows()

    return True


def check_area(contours, area_thres_min=1200, area_thres_max=50000) -> list:
    valid_contours = []
    for i, contour in enumerate(contours):
        x, y, w, h = cv.boundingRect(contour)
        print('area : ', w * h)
        if w / h > 6 or w / h < 0.15:
            continue
        if w * h < area_thres_min or w * h > area_thres_max:
            continue
        # if w / h > 3.5:
        #     continue
        valid_contours.append(contour)

    return valid_contours


def remove_inner_box(contours):
    outer_contours = []
    for contour1 in contours:
        is_inner_box = False
        x1_min, y1_min, w1, h1 = cv.boundingRect(contour1)
        x1_max = x1_min + w1
        y1_max = y1_min + h1
        for contour2 in contours:
            x2_min, y2_min, w2, h2 = cv.boundingRect(contour2)
            x2_max = x2_min + w2
            y2_max = y2_min + h2
            if x1_min > x2_min and x1_max < x2_max and y1_min > y2_min and y1_max < y2_max:
                is_inner_box = True
                break
        if not is_inner_box:
            outer_contours.append(contour1)

    return outer_contours


def sort_contour(contours):
    sorted_contours = []
    # x축 방향 정렬
    contours = sorted(contours, key=lambda contour: cv.boundingRect(contour)[0])

    return sort_contour_y(contours, sorted_contours)


def sort_contour_y(contours, sorted_contours):
    for contour1 in contours:
        if is_contain_contour(sorted_contours, contour1):
            continue
        is_next_box = True
        for contour2 in contours:
            if is_contain_contour(sorted_contours, contour2):
                continue
            if not get_is_next_box(contour1, contour2, contours):
                is_next_box = False
                break
        if not is_next_box:
            continue
        sorted_contours.append(contour1)
        break

    if len(contours) == len(sorted_contours):
        return sorted_contours

    return sort_contour_y(contours, sorted_contours)


def is_contain_contour(sorted_contours, contour) -> bool:
    box1 = cv.boundingRect(contour)
    for temp in sorted_contours:
        box2 = cv.boundingRect(temp)
        if box1[0] == box2[0] and box1[1] == box2[1] and box1[2] == box2[2] and box1[3] == box2[3]:
            return True
    return False


def get_is_next_box(contour1, contour2, contours) -> bool:
    x1, y1, w1, h1 = cv.boundingRect(contour1)
    x2, y2, w2, h2 = cv.boundingRect(contour2)
    if y1 + h1 / 3 >= y2 + h2:
        # 같은 라인(row)에 있는 box 인지 확인
        for contour in contours:
            if is_overlap_y(contour1, contour) and is_overlap_y(contour2, contour) and not is_overlap_x(contour1, contour2):
                return True
        return False
    if y1 + h1 <= y2:
        return True
    if x1 + w1 <= x2:
        return True

    return True


def merge_box(boxes):
    return boxes


def merge(a, b):
    x = min(a[0], b[0])
    y = min(a[1], b[1])
    w = max(a[0] + a[2], b[0] + b[2]) - x
    h = max(a[1] + a[3], b[1] + b[3]) - y
    return [x, y, w, h]


def is_overlap_y(contour1, contour2) -> bool:
    x1, y1, w1, h1 = cv.boundingRect(contour1)
    x2, y2, w2, h2 = cv.boundingRect(contour2)

    if y2 <= y1 <= y2 + h2 or y2 <= y1 + h1 <= y2 + h2:
        return True

    if y1 <= y2 <= y1 + h1 or y1 <= y2 + h2 <= y1 + h1:
        return True

    return False


def is_overlap_x(contour1, contour2) -> bool:
    x1, y1, w1, h1 = cv.boundingRect(contour1)
    x2, y2, w2, h2 = cv.boundingRect(contour2)

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


def is_valid_contour(index, contour, contours, src) -> bool:
    # approx = cv.approxPolyDP(contour, 10, True)
    # cv.drawContours(src, approx, 0, (0, 0, 255), thickness=3)
    # cv.rectangle(src, cv.boundingRect(approx), (0, 0, 255), thickness=1)
    # if not cv.isContourConvex(approx):
    #     return False

    new_height, new_width = src.shape[:2]
    x, y, w, h = cv.boundingRect(contour)
    if w / h > 6 or w / h < 0.15:
        return False
    if w * h < 1200:  # or w * h > 50000:
        return False
    if w / new_width > 0.9 or h / new_height > 0.9:
        return False
    return True


def remove_invalid_box(boxes):
    valid_boxes = []
    if len(boxes) == 0:
        return valid_boxes

    area = boxes[0][2] * boxes[0][3]
    if area < 2500:
        boxes.pop(0)
    area = boxes[len(boxes) - 1][2] * boxes[len(boxes) - 1][3]
    if area < 2500:  # or area > 70000:
        boxes.pop(len(boxes) - 1)

    box_areas = []
    # valid_boxes = []
    for box in boxes:
        area = box[2] * box[3]
        box_areas.append(area)
        print('area : ', area)

    # box_areas.remove(max(box_areas))
    # box_areas.remove(min(box_areas))

    average = np.mean(box_areas)
    std = np.std(box_areas)
    print('average : ', average)
    print('std     : ', std)
    valid_boxes = filter(lambda b: abs((b[2] * b[3] - average) / std) < 4.0, boxes)

    return valid_boxes


if __name__ == "__main__":
    main()
    # crop_text('')
