FROM quay.io/nordstrom/baseimage-ubuntu:16.04

ADD gettomethod-linux-amd64 /gettomethod

ENTRYPOINT /gettomethod
