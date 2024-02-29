import os
import sys
import unicodedata
import numpy as np

import cv2 as cv

from ch07.vehicle_no_2 import ChainParam, DetectCustom, DetectYolo
import easyocr
from hangul_jamo import compose

train_root = '/Users/kcson/mywork/data/lpr_train'


def main():
    # root = '/Users/kcson/mywork/data/[원천]자동차번호판OCR데이터'
    # root = '../imgs/car'
    # root = '/Users/kcson/mywork/data/lpr_pre'
    root = '/Users/kcson/mywork/data/lpr_pre_auto_gen'
    train_file_list = os.listdir(root)

    index = 0
    valid_index = 0
    h_sum = 0
    w_sum = 0
    for file_name in train_file_list:
        # if index == 150:
        #     break
        file_name = unicodedata.normalize('NFC', file_name)
        label = os.path.splitext(file_name)[0].split('-')[0]
        label = eng_to_region(label)
        label = eng_to_kor(label)
        label = unicodedata.normalize('NFC', label)

        print(label)
        # if len(label) > 8:
        #     continue

        full_path = os.path.join(root, file_name)
        # if os.path.getsize(full_path) < 30 * 1024:
        #     continue
        print(full_path)
        try:
            if crop_text(full_path, label):
                index += 1
                print('index : ', index)
        except Exception as e:
            print(e)

    print('index : ', index)


def crop_text(full_path, label=''):
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

    detect = DetectYolo()
    detect.set_next(DetectCustom())
    # detect = DetectCustom()
    # detect.set_next(DetectYolo())
    detect.handle(chainParam)

    chainParam.dst = src
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
    src_bin_g = cv.GaussianBlur(src_bin, (0, 0), 3, sigmaY=0, borderType=cv.BORDER_DEFAULT)
    # src_bin = cv.adaptiveThreshold(src_bin, 255, cv.ADAPTIVE_THRESH_GAUSSIAN_C, cv.THRESH_BINARY_INV, 19, 4)
    thres, _ = cv.threshold(src_bin_g, 0, 255, cv.THRESH_BINARY | cv.THRESH_OTSU, dst=src_bin)
    kernel = cv.getStructuringElement(cv.MORPH_RECT, (6, 6))
    # src_bin = cv.dilate(src_bin, kernel, iterations=1)
    # src_bin = cv.erode(src_bin, kernel, iterations=2)

    # src_bin = cv.inRange(src, (150, 150, 150), (255, 255, 255))
    # thres, _ = cv.threshold(src_bin, 0, 255, cv.THRESH_BINARY_INV | cv.THRESH_OTSU, dst=src_bin)

    cv.imshow('src_bin', src_bin)
    cv.waitKey()

    contours, _ = cv.findContours(src_bin, cv.RETR_LIST, cv.CHAIN_APPROX_NONE)

    # contour -> box
    boxes = contours_to_boxes(contours)

    # contour 면적이 너무 작거나 큰 box 제거
    boxes = check_area(boxes, area_thres_min=750, area_thres_max=60000, ratio_thres_min=0.11, ratio_thres_max=5.5)

    # 내부에 있는 box 제거
    boxes = remove_inner_box(boxes)

    # 면적이 다른 box 와 차이가 많이 나는 box 제거
    boxes = check_area_diff(boxes, min_thres=0.1, max_thres=3.5)

    # y 축으로 떨어져 있는 box 제거
    boxes = check_y_diff(boxes)

    # box 정렬(위->아래, 좌->우)
    boxes = sort_box(boxes)

    # 기타 유효 하지 않은 box 제거
    boxes = remove_invalid_box(boxes, min_area=2000)

    # 면적이 다른 box 와 차이가 많이 나는 box 제거
    boxes = check_area_diff(boxes, min_thres=0.1, max_thres=3.5)

    # box merge
    boxes, row_index = merge_box(boxes)

    # 기타 유효 하지 않은 box 제거
    boxes = remove_invalid_box(boxes, min_area=2500)

    # contour 면적이 너무 작거나 큰 box 제거
    boxes = check_area(boxes, area_thres_min=2000, area_thres_max=60000, ratio_thres_max=2.0)

    # y 축으로 떨어져 있는 box 제거
    boxes = check_y_diff(boxes)

    # 면적이 다른 box 와 차이가 많이 나는 box 제거
    boxes = check_area_diff(boxes, min_thres=0.5, max_thres=2.0)

    if row_index == 1:
        thres, _ = cv.threshold(src_bin_g, 0, 255, cv.THRESH_BINARY | cv.THRESH_OTSU, dst=src_bin)

    # boxes.pop(0)
    # boxes.pop(4)
    # temp = merge(boxes[2],boxes[3])
    # boxes[2] = temp
    # boxes.pop(3)
    for i, box in enumerate(boxes, start=1):
        x, y, w, h = box
        print('area : ', h * w, ' ratio : ', w / h)
        cv.rectangle(src, (x, y, w, h), (0, 0, 255), thickness=1)
        cv.putText(src, str(i), (x, y), cv.FONT_HERSHEY_SIMPLEX, 1, (0, 0, 255), 2, cv.LINE_AA)

        if i > len(label)-1:
            continue
        train_image = src_bin[y:y + h, x:x + w]
        train_image = cv.resize(train_image, (80, 160))
        train_path = os.path.join(train_root, label[i - 1])
        try:
            file_count = len(os.listdir(train_path))
        except:
            file_count = 0
            os.mkdir(train_path)

        train_path = os.path.join(train_path, '{}.jpeg'.format(file_count + 1))
        print(train_path, " : ", file_count)
        cv.imwrite(train_path, train_image)

    # recognition_text(boxes, src_bin)
    cv.imshow('src_bin', src_bin)
    cv.imshow('src', src)
    cv.waitKey()
    cv.destroyAllWindows()

    if label != '' and len(label) != len(boxes):
        return False

    return True


def check_y_diff(boxes, valid_thres=3):
    valid_boxes = []
    for box1 in boxes:
        valid_count = 0
        for box2 in boxes:
            if is_overlap_y(box1, box2):
                valid_count += 1
        if valid_count < valid_thres:
            continue
        valid_boxes.append(box1)

    return valid_boxes


def check_area_diff(boxes, valid_thres=3, min_thres=0.5, max_thres=2.0):
    valid_boxes = []
    for box1 in boxes:
        area1 = box1[2] * box1[3]
        valid_count = 0
        for box2 in boxes:
            area2 = box2[2] * box2[3]
            if min_thres < area2 / area1 < max_thres:
                valid_count += 1
        if valid_count < valid_thres:
            continue
        valid_boxes.append(box1)

    return valid_boxes


def contours_to_boxes(contours):
    boxes = []
    for contour in contours:
        boxes.append(cv.boundingRect(contour))

    return boxes


def check_area(boxes, area_thres_min=1200, area_thres_max=50000, ratio_thres_min=0.15, ratio_thres_max=6.0) -> list:
    valid_boxes = []
    for i, box in enumerate(boxes):
        x, y, w, h = box
        print('area : ', w * h)
        if w / h > ratio_thres_max or w / h < ratio_thres_min:
            continue
        if w * h < area_thres_min or w * h > area_thres_max:
            continue
        # if w / h > 3.5:
        #     continue
        valid_boxes.append(box)

    return valid_boxes


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


def sort_box(boxes):
    sorted_boxes = []
    # x축 방향 정렬
    boxes = sorted(boxes, key=lambda box: box[0])

    return sort_box_y(boxes, sorted_boxes)


def sort_box_y(boxes, sorted_boxes):
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

    return sort_box_y(boxes, sorted_boxes)


def is_contain_contour(sorted_contours, contour) -> bool:
    box1 = cv.boundingRect(contour)
    for temp in sorted_contours:
        box2 = cv.boundingRect(temp)
        if box1[0] == box2[0] and box1[1] == box2[1] and box1[2] == box2[2] and box1[3] == box2[3]:
            return True
    return False


def get_is_next_box(box1, box2, boxes) -> bool:
    x1, y1, w1, h1 = box1
    x2, y2, w2, h2 = box2
    if y1 + h1 / 3 >= y2 + h2:
        # 같은 라인(row)에 있는 box 인지 확인
        for box in boxes:
            if is_overlap_y(box1, box) and is_overlap_y(box2, box) and not is_overlap_x(box1, box2):
                return True
        return False
    if y1 + h1 <= y2:
        return True
    if x1 + w1 <= x2:
        return True

    return True


def merge_box(boxes):
    return_row_index = 1
    index = 0
    merged_boxes = []
    while index < len(boxes):
        merge_box_count, row_index = get_merge_box_count(index, boxes)
        if row_index == 0:
            return_row_index = 0
            start_box = boxes[index]
            for i in range(index + 1, index + merge_box_count):
                start_box = merge(start_box, boxes[i])
            merged_boxes.append(start_box)
            index += merge_box_count
        else:
            if merge_box_count == 1:
                merged_boxes.append(boxes[index])
                index += 1
            elif merge_box_count == 2:
                merged_boxes.append(boxes[index])
                merged_boxes.append(boxes[index + 1])
                index += 2
            elif merge_box_count == 3:
                merged_boxes.append(boxes[index])
                merged_boxes.append(merge(boxes[index + 1], boxes[index + 2]))
                index += 3
            elif merge_box_count == 4:
                merged_boxes.append(merge(boxes[index], boxes[index + 1]))
                merged_boxes.append(merge(boxes[index + 2], boxes[index + 3]))
                index += 4
            elif merge_box_count == 5:
                merged_boxes.append(merge(merge(boxes[index], boxes[index + 1]), boxes[index + 2]))
                merged_boxes.append(merge(boxes[index + 3], boxes[index + 4]))
                index += 5
            elif merge_box_count == 6:
                merged_boxes.append(merge(merge(boxes[index], boxes[index + 1]), boxes[index + 2]))
                merged_boxes.append(merge(merge(boxes[index + 3], boxes[index + 4]), boxes[index + 5]))
                index += 6
            else:
                merged_boxes.append(boxes[index])
                index += 1

    index = 0
    result_boxes = []
    while index < len(merged_boxes):
        box1 = merged_boxes[index]
        if index == len(merged_boxes) - 1:
            result_boxes.append(box1)
            break
        box2 = merged_boxes[index + 1]
        if is_overlap_x(box1, box2) and index > 0:
            result_boxes.append(merge(box1, box2))
            index += 2
            continue
        if box2[2] / box2[3] < 0.35:
            if (index + 1 == len(merged_boxes) - 5) or \
                    (get_row_index(box1, boxes) == 0 and index + 1 == 2):
                result_boxes.append(merge(box1, box2))
                index += 2
                continue
        result_boxes.append(box1)
        index += 1
    return result_boxes, return_row_index


def get_merge_box_count(index, boxes):
    row_index = 0
    max_index = index
    start_boxs = [boxes[index]]

    for i in range(index + 1, len(boxes)):
        box1 = boxes[i]
        for start_box in start_boxs:
            overlap_x = False
            if is_overlap_x(start_box, box1):
                for box in boxes:
                    if is_overlap_y(start_box, box) and is_overlap_y(box1, box):
                        max_index = i
                        start_boxs.append(box1)
                        overlap_x = True
                        break
                if overlap_x:
                    break
            if box1[2] / box1[3] < 0.2 and get_row_index(box1,boxes) == 0:
                max_index = i
                break

    row_index = get_row_index(start_boxs[len(start_boxs) - 1], boxes)

    merge_box_count = max_index - index + 1
    return merge_box_count, row_index


def get_row_index(box, boxes):
    under_box_count = 0
    for box1 in boxes:
        if box[1] + box[3] < box1[1]:
            for box2 in boxes:
                if is_overlap_y(box, box2) and is_overlap_y(box1, box2):
                    continue
            under_box_count += 1
    if under_box_count < 4:
        return 1
    return 0


def merge(a, b):
    x = min(a[0], b[0])
    y = min(a[1], b[1])
    w = max(a[0] + a[2], b[0] + b[2]) - x
    h = max(a[1] + a[3], b[1] + b[3]) - y
    return [x, y, w, h]


def is_overlap_y(box1, box2, epsilon=0) -> bool:
    x1, y1, w1, h1 = box1
    x2, y2, w2, h2 = box2

    if y2 <= y1 - epsilon <= y2 + h2 or y2 <= y1 + h1 - epsilon <= y2 + h2:
        return True

    if y1 <= y2 - epsilon <= y1 + h1 or y1 <= y2 + h2 - epsilon <= y1 + h1:
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


def remove_invalid_box(boxes, min_area=2500, max_area=70000, min_ratio=0.95):
    valid_boxes = []
    if len(boxes) == 0:
        return valid_boxes

    area = boxes[0][2] * boxes[0][3]
    if area < min_area:
        boxes.pop(0)
    area = boxes[len(boxes) - 1][2] * boxes[len(boxes) - 1][3]
    ratio = boxes[len(boxes) - 1][2] / boxes[len(boxes) - 1][3]
    if area < min_area or ratio > min_ratio:  # or area > 70000:
        boxes.pop(len(boxes) - 1)

    for box in boxes:
        valid_boxes.append(box)

    return valid_boxes


# 자음-초성/종성
cons = {'r': 'ㄱ', 'R': 'ㄲ', 's': 'ㄴ', 'e': 'ㄷ', 'E': 'ㄸ', 'f': 'ㄹ', 'a': 'ㅁ', 'q': 'ㅂ', 'Q': 'ㅃ', 't': 'ㅅ', 'T': 'ㅆ',
        'd': 'ㅇ', 'w': 'ㅈ', 'W': 'ㅉ', 'c': 'ㅊ', 'z': 'ㅋ', 'x': 'ㅌ', 'v': 'ㅍ', 'g': 'ㅎ'}
# 모음-중성
vowels = {'k': 'ㅏ', 'o': 'ㅐ', 'i': 'ㅑ', 'O': 'ㅒ', 'j': 'ㅓ', 'p': 'ㅔ', 'u': 'ㅕ', 'P': 'ㅖ', 'h': 'ㅗ', 'hk': 'ㅘ', 'ho': 'ㅙ', 'hl': 'ㅚ',
          'y': 'ㅛ', 'n': 'ㅜ', 'nj': 'ㅝ', 'np': 'ㅞ', 'nl': 'ㅟ', 'b': 'ㅠ', 'm': 'ㅡ', 'ml': 'ㅢ', 'l': 'ㅣ'}
# 자음-종성
cons_double = {'rt': 'ㄳ', 'sw': 'ㄵ', 'sg': 'ㄶ', 'fr': 'ㄺ', 'fa': 'ㄻ', 'fq': 'ㄼ', 'ft': 'ㄽ', 'fx': 'ㄾ', 'fv': 'ㄿ', 'fg': 'ㅀ', 'qt': 'ㅄ'}


# A = 서울 B = 경기 C = 인천 D = 강원 E = 충남, F = 대전 G = 충북 H = 부산 I = 울산 J =대구 K = 경북 L = 경남 M = 전남 N = 광주 O = 전북 P = 제주
def eng_to_region(text):
    text = text.replace('A', '서울'). \
        replace('B', '경기'). \
        replace('C', '인천'). \
        replace('D', '강원'). \
        replace('E', '충남'). \
        replace('F', '대전'). \
        replace('G', '충북'). \
        replace('H', '부산'). \
        replace('I', '울산'). \
        replace('J', '대구'). \
        replace('K', '경북'). \
        replace('L', '경남'). \
        replace('M', '전남'). \
        replace('N', '광주'). \
        replace('O', '전북'). \
        replace('P', '제주')

    text = text.replace('Z', '').replace('X', '')

    return text


def eng_to_kor(text):
    result = ''  # 영 > 한 변환 결과

    # 1. 해당 글자가 자음인지 모음인지 확인
    vc = ''
    for t in text:
        if t in cons:
            vc += 'c'
        elif t in vowels:
            vc += 'v'
        else:
            vc += '!'

    # cvv → fVV / cv → fv / cc → dd
    vc = vc.replace('cvv', 'fVV').replace('cv', 'fv').replace('cc', 'dd')

    # 2. 자음 / 모음 / 두글자 자음 에서 검색
    i = 0
    while i < len(text):
        v = vc[i]
        t = text[i]

        j = 1
        # 한글일 경우
        try:
            if v == 'f' or v == 'c':  # 초성(f) & 자음(c) = 자음
                result += cons[t]

            elif v == 'V':  # 더블 모음
                result += vowels[text[i:i + 2]]
                j += 1

            elif v == 'v':  # 모음
                result += vowels[t]

            elif v == 'd':  # 더블 자음
                result += cons_double[text[i:i + 2]]
                j += 1
            else:
                result += t

        # 한글이 아닐 경우
        except:
            if v in cons:
                result += cons[t]
            elif v in vowels:
                result += vowels[t]
            else:
                result += t

        i += j

    return compose(result)


if __name__ == "__main__":
    # main()
    crop_text('/Users/kcson/mywork/data/lpr_sample/vehicle45.jpeg', '경기38더5484')
