from langchain.schema import BaseMemory
from pydantic import BaseModel
from typing import List, Dict, Any

class Bridge:
    def open_index_session(self, index: str, namespace: str) -> "BridgeIndexSession":
        pass

class BridgeIndexSession:
    bridge: Bridge

    index: str
    namespace: str

    root_memory_id: str
    branch_memory_id: str
    parent_memory_id: str
    current_memory_id: str

    current_clock: int
    current_height: int

    def get_memory_data(self) -> Dict[str, any]:
        pass

    def update_memory_data(self, data: Dict[str, any]):
        pass

    def discard(self):
        pass

    def merge(self):
        pass

class BridgeMemory(BaseMemory, BaseModel):
    bridge: Bridge
    session: BridgeIndexSession
    auto_commit: bool = False
    memory_key: str = "memory"

    def clear(self):
        pass

    @property
    def memory_variables(self) -> List[str]:
        return [self.memory_key]

    def load_memory_variables(self, inputs: Dict[str, Any]) -> Dict[str, str]:
        context = self.session.get_memory_data()

        return {self.memory_key: context}

    def save_context(self, inputs: Dict[str, Any], outputs: Dict[str, str]) -> None:
        self.session.update_memory_data({
            "inputs": inputs,
            "outputs": outputs,
        })

        if self.auto_commit:
            self.session.merge()
