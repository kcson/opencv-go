import os
import sys
from random import randint

import cv2 as cv

rotate_degrees = [-7, -6, -5, -4, 0, 0, 0, 0, 0, 0, 0, 4, 5, 6, 7]
root = '/Users/kcson/mywork/data/lpr_train_final'
root_dir = os.listdir(root)
# for sub_dir in root_dir:
#     if sub_dir.startswith('.'):
#         continue
# sub_dir = os.path.join(root, sub_dir)
sub_dir = os.path.join(root, '울-1')
file_index = 1
augment_file_index = 1001
while file_index <= 500:
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

    cut_length = randint(5, 10)
    cut_direction = randint(1, 4)
    if cut_direction == 1:
        cut_image = cv.copyMakeBorder(src, cut_length, 0, 0, 0, cv.BORDER_CONSTANT, value=(255, 255, 255))
        cut_image = cut_image[:src_h, :]
    elif cut_direction == 2:
        cut_image = cv.copyMakeBorder(src, 0, cut_length, 0, 0, cv.BORDER_CONSTANT, value=(255, 255, 255))
        cut_image = cut_image[cut_length:, :]
    elif cut_direction == 3:
        cut_image = cv.copyMakeBorder(src, 0, 0, cut_length, 0, cv.BORDER_CONSTANT, value=(255, 255, 255))
        cut_image = cut_image[:, :src_w]
    else:
        cut_image = cv.copyMakeBorder(src, 0, 0, 0, cut_length, cv.BORDER_CONSTANT, value=(255, 255, 255))
        cut_image = cut_image[:, cut_length:]

    dilate_iter = randint(1, 3)
    dilate_image = cv.dilate(cut_image, None, iterations=dilate_iter)

    # cv.imshow('cut_image', cut_image)
    # cv.waitKey()
    # cv.imshow('dilate_image', dilate_image)
    # cv.waitKey()

    cut_file = os.path.join(sub_dir, '{}.jpeg'.format(augment_file_index))
    cv.imwrite(cut_file, cut_image)
    augment_file_index += 1
    dilate_file = os.path.join(sub_dir, '{}.jpeg'.format(augment_file_index))
    cv.imwrite(dilate_file, dilate_image)
    augment_file_index += 1
    # cv.imwrite(target_file, src_rotated)

# cv.imshow('src', src)
# cv.imshow('src_padding', src_padding)
# cv.imshow('src_rotated', src_rotated)
# cv.waitKey()
# cv.destroyAllWindows()
