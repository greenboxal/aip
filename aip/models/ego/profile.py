from typing import List

from python_easy_json import JSONObject

class BaseModel(JSONObject):
    pass

class Aptitude(BaseModel):
    description: str
    self_reflection: str

class Desire(BaseModel):
    description: str
    self_reflection: str


class Goal(BaseModel):
    description: str
    self_reflection: str

class Profile(BaseModel):
    name: str
    directive: str

    executive_model: str

    aptitudes: List[Aptitude]
    desires: List[Desire]
    goals: List[Goal]

    def __init__(self, **kwargs):
        self.aptitudes = []
        self.desires = []
        self.goals = []

        super().__init__(**kwargs)

class MindState(BaseModel):
    profile: Profile
    aptitudes: [Aptitude]
    desires: [Desire]
    goals: [Goal]

    description: str
    self_reflection: str

    def __init__(self, **kwargs):
        self.aptitudes = []
        self.desires = []
        self.goals = []

        super().__init__(**kwargs)
