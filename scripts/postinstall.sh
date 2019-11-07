#!/bin/bash
systemctl daemon-reload
systemctl enable hkick-core
systemctl start hkick-core