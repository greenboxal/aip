from langchain.schema import BaseMemory
from pydantic import BaseModel
from typing import List, Dict, Any
import requests

class BridgeMemory(BaseMemory, BaseModel):
    memory_key: str = "memory"
    endpoint: str = "http://localhost:30100/v1/rpc"

    _current_memory: any

    def clear(self):
        pass

    @property
    def memory_variables(self) -> List[str]:
        return [self.memory_key]

    def load_memory_variables(self, inputs: Dict[str, Any]) -> Dict[str, str]:
        memory_id = None

        if self._current_memory is not None:
            memory_id = self._current_memory.metadata.id

        result = self._send_request("memlink.OneShotGetMemory", {
            "memory_id": memory_id,
        })

        self._current_memory = result.memory

        return {self.memory_key: self._current_memory.data.text}

    def save_context(self, inputs: Dict[str, Any], outputs: Dict[str, str]) -> None:
        filtered_inputs = {k: v for k, v in inputs.items() if k != self.memory_key}

        texts = [
            f"{k}: {v}"
            for k, v in list(filtered_inputs.items()) + list(outputs.items())
        ]

        page_content = "\n".join(texts)

        self._current_memory = self._send_request("memlink.OneShotPutMemory", {
            "old_memory": self._current_memory,
            "new_memory": {
                "text": page_content,
            }
        })

    def _send_request(self, method, request):
        payload = {
            "method": method,
            "params": [request],
            "jsonrpc": "2.0",
            "id": 0,
        }

        return requests.post(self.endpoint, json=payload).json()
