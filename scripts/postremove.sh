#!/bin/bash
systemctl disable kickr-core
systemctl stop kickr-core
systemctl daemon-reload