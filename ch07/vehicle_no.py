from __future__ import annotations
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
        print('PreProcess')
        cv.imshow('src', param.src)
        return super().handle(param)


class CutVehicleRegion(AbstractHandler):
    def handle(self, param: ChainParam):
        param.dst = param.src
        print('CutVehicleRegion')
        return super().handle(param)


class GetVehicleNo(AbstractHandler):
    def handle(self, param: ChainParam):
        print('GetVehicleNo')
        cv.imshow('dst==', param.dst)
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

    # cv.imshow('src', src)

    cv.waitKey()
    cv.destroyAllWindows()
