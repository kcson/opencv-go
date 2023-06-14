import os
import unicodedata
import cv2 as cv

root = '/Users/kcson/mywork/data/[원천]자동차번호판OCR데이터'

train_file_list = os.listdir(root)

for file_name in train_file_list:
    file_name = unicodedata.normalize('NFC', file_name)
    label = os.path.splitext(file_name)[0].split('-')[0]
    print(label, " : ", len(label))

