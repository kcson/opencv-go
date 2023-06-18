import math
import os
import unicodedata
import cv2 as cv

from ch07.vehicle_no_2 import ChainParam, DetectCustom, DetectYolo


def main():
    root = '/Users/kcson/mywork/data/[원천]자동차번호판OCR데이터'

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

        fullPath = os.path.join(root, file_name)

        src = cv.imread(fullPath)
        if src is None:
            continue
        height, width = src.shape[:2]
        if width < 300:
            continue

        chainParam = ChainParam(src)
        chainParam.copySrc = src.copy()

        detect = DetectCustom()
        detect.handle(chainParam)
        # chainParam.dst = src.copy()
        if chainParam.dst is not None:
            height, width = chainParam.dst.shape[:2]
            ratio = width / height
            new_width = 1024
            new_height = int(new_width / ratio)
            chainParam.dst = cv.resize(chainParam.dst, (new_width, new_height))

            src_bin = cv.cvtColor(chainParam.dst, cv.COLOR_BGR2GRAY)
            src_bin = cv.GaussianBlur(src_bin, (0, 0), 2, sigmaY=0, borderType=cv.BORDER_DEFAULT)
            # src_bin = cv.bilateralFilter(src_bin, 0, 10, 5)
            thres, _ = cv.threshold(src_bin, 0, 255, cv.THRESH_BINARY | cv.THRESH_OTSU, dst=src_bin)

            src_bin = cv.erode(src_bin, None, iterations=6)
            # src_bin = cv.dilate(src_bin, None, iterations=1)

            contours, _ = cv.findContours(src_bin, cv.RETR_LIST, cv.CHAIN_APPROX_NONE)

            train_image_count = 0
            for contour in contours:
                x, y, w, h = cv.boundingRect(contour)
                if w / h > 1.9:
                    continue
                if w * h < 2000:
                    continue
                # if h < 100:
                #     continue
                # if similar_contour(contours, contour) < 7:
                #     continue
                print(w, h, w * h)
                h_sum += h
                w_sum += w
                valid_index += 1
                train_image_count += 1
                cv.rectangle(chainParam.dst, (x, y, w, h), (0, 0, 255), thickness=1)
                cv.putText(chainParam.dst, str(train_image_count), (x, y), cv.FONT_HERSHEY_SIMPLEX, 2, (0, 0, 255), 1, cv.LINE_AA)

            # if train_image_count != len(label):
            #     continue

            cv.imshow('src', src)
            cv.imshow('src_bin', src_bin)
            cv.imshow('dst', chainParam.dst)
            cv.waitKey()

            print(label, " : ", len(label))
            index += 1

    print(w_sum / valid_index, h_sum / valid_index, index)
    cv.destroyAllWindows()


def similar_contour(contours, pts) -> int:
    similar_count = 0
    x1, y1, w1, h1 = cv.boundingRect(pts)
    diagonal_length = math.sqrt(w1 * w1 + h1 * h1)
    for pts in contours:
        x2, y2, w2, h2 = cv.boundingRect(pts)
        area_diff: float = abs(w1 * h1 - w2 * h2) / (w1 * h1)
        width_diff: float = abs(w1 - w2) / w1
        height_diff: float = abs(h1 - h2) / h1

        if area_diff > 0.9 or width_diff > 0.8:# or height_diff > 0.7:
            continue

        # distance = math.sqrt((x1 - x2) ** 2 + (y1 - y2) ** 2)
        # if distance > diagonal_length * 10:
        #     continue
        #
        # angle_diff = 0
        # dx = abs(x1 - x2)
        # dy = abs(y1 - y2)
        # if dx == 0:
        #     angle_diff = 90
        # else:
        #     angle_diff = math.atan(dy / dx) * (180 / math.pi)
        # if angle_diff > 20.0:
        #     continue

        similar_count = similar_count + 1
    return similar_count


if __name__ == "__main__":
    main()
