import os
import sys
from random import randint

import cv2 as cv

rotate_degrees = [-7, -6, -5, -4, 0, 0, 0, 0, 0, 0, 0, 4, 5, 6, 7]
root = '/Users/kcson/mywork/data/lpr_train_final'
root_dir = os.listdir(root)
for sub_dir in root_dir:
    if sub_dir.startswith('.'):
        continue
    sub_dir = os.path.join(root, sub_dir)
    file_index = 2
    while file_index <= 1000:
        image_file = os.path.join(sub_dir, '{}.jpeg'.format(file_index))
        target_file = image_file
        file_index += 1
        if not os.path.exists(image_file):
            image_file = os.path.join(sub_dir, '1.jpeg')
        print(target_file, ' : ', image_file)

        src = cv.imread(image_file)
        if src is None:
            print("image read fail!!")
            sys.exit()

        src = cv.cvtColor(src, cv.COLOR_BGR2GRAY)
        src_h, src_w = src.shape

        src_padding = cv.copyMakeBorder(src, 10, 10, 10, 10, cv.BORDER_CONSTANT, value=(255, 255, 255))
        src_padding_h, src_padding_w = src_padding.shape

        rotate_figure = cv.getRotationMatrix2D((int(src_padding_w / 2), int(src_padding_h / 2)), rotate_degrees[randint(0, 14)], 1)
        src_rotated = cv.warpAffine(src_padding, rotate_figure, (src_padding_w, src_padding_h), borderValue=(255, 255, 255))

        src_rotated = cv.resize(src_rotated, (src_w, src_h))
        cv.imwrite(target_file, src_rotated)

# cv.imshow('src', src)
# cv.imshow('src_padding', src_padding)
# cv.imshow('src_rotated', src_rotated)
# cv.waitKey()
# cv.destroyAllWindows()
