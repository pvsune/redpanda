# Copyright 2020 Vectorized, Inc.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.md
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the Apache License, Version 2.0

import sys
from ducktape.services.background_thread import BackgroundThreadService
from ducktape.cluster.remoteaccount import RemoteCommandError
from threading import Event


class KafProducer(BackgroundThreadService):
    def __init__(self, context, redpanda, topic, num_records=sys.maxsize):
        super(KafProducer, self).__init__(context, num_nodes=1)
        self._redpanda = redpanda
        self._topic = topic
        self._num_records = num_records
        self._stopping = Event()

    def _worker(self, _idx, node):
        cmd = f"for (( i=0; i < {self._num_records}; i++ )) ; do export KEY=key-$(printf %08d $i) ; export VALUE=record-$(printf %08d $i) ; echo $VALUE | kaf produce -b {self._redpanda.brokers()} --key $KEY {self._topic} ; done"

        self._stopping.clear()
        try:
            for line in node.account.ssh_capture(cmd, timeout_sec=10):
                self.logger.debug(line.rstrip())
        except RemoteCommandError:
            if self._stopping.is_set():
                pass
            else:
                raise

    def stop_node(self, node):
        self._stopping.set()
        try:
            node.account.kill_process("kaf", clean_shutdown=False)
        except RemoteCommandError as e:
            if b"No such process" in e.msg:
                pass
            else:
                raise

    def clean_node(self, nodes):
        pass