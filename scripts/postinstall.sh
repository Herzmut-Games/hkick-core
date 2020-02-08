#!/bin/bash
systemctl daemon-reload
systemctl enable kickr-core
systemctl start kickr-core