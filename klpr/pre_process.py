import os
import unicodedata
import cv2 as cv

from ch07.vehicle_no_2 import ChainParam, DetectCustom

root = '/Users/kcson/mywork/data/[원천]자동차번호판OCR데이터'

train_file_list = os.listdir(root)

index = 0
for file_name in train_file_list:
    if index == 10:
        break
    file_name = unicodedata.normalize('NFC', file_name)
    label = os.path.splitext(file_name)[0].split('-')[0]
    if len(label) > 8:
        continue

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
    if chainParam.dst is not None:
        cv.imshow('dst', chainParam.dst)
        cv.imshow('src', src)
        cv.waitKey()

        print(label, " : ", len(label))
        index += 1

cv.destroyAllWindows()
