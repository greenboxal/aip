from typing import List

from langchain import OpenAI, WolframAlphaAPIWrapper, ArxivAPIWrapper
from langchain.chat_models.base import BaseChatModel
from langchain.llms import BaseLLM
from langchain.utilities import PythonREPL
from langchain.vectorstores.base import VectorStoreRetriever

from aip.models.ego.profile import Profile, MindState
from langchain.memory.chat_memory import BaseMemory
from langchain.chains import ConversationChain
from langchain.chat_models import ChatOpenAI
from langchain.prompts.prompt import PromptTemplate
from langchain.agents import initialize_agent
from langchain.agents import AgentExecutor, AgentType, Tool
from langchain.tools.base import BaseTool
from langchain.agents import load_tools

from langchain.prompts import (
    ChatPromptTemplate,
    ChatMessagePromptTemplate,
    MessagesPlaceholder,
    HumanMessagePromptTemplate,
    AIMessagePromptTemplate,
    SystemMessagePromptTemplate,
)

from langchain.memory import (
    ChatMessageHistory,
    ConversationSummaryBufferMemory,
    CombinedMemory,
    VectorStoreRetrieverMemory,
    ReadOnlySharedMemory,
)




class Persona:
    profile: Profile
    state: MindState

    llm: BaseLLM
    chat_llm: BaseChatModel
    memory: BaseMemory
    monologue: ConversationChain
    tools: List[BaseTool]
    agent: AgentExecutor

    def __init__(self, profile: Profile, retriever: VectorStoreRetriever, verbose=True):
        self.profile = profile
        self.verbose = verbose

        self.state = MindState(data={"profile": self.profile})

        self.llm = ChatOpenAI(temperature=0)
        self.chat_llm = ChatOpenAI(temperature=0.5, model_name="gpt-3.5-turbo")

        self.vec_memory = VectorStoreRetrieverMemory(
            memory_key="context",
            input_key="input",
            retriever=retriever,
            return_docs=True,
        )

        self.readonly_vec_memory = ReadOnlySharedMemory(memory=self.vec_memory)

        self.short_term_memory = ConversationSummaryBufferMemory(
            llm=self.llm,
            max_token_limit=200,
            memory_key="chat_history",
            input_key="input",
            ai_prefix=self.profile.name,
        )

        self.memory = CombinedMemory(
            memories=[
                self.short_term_memory,
                self.readonly_vec_memory
            ],
        )

        prompt = self._build_monologue_prompt()

        self.monologue = ConversationChain(
            llm=self.chat_llm,
            memory=self.memory,
            prompt=prompt,
            verbose=verbose,
        )

        self.tools = []

        self.tools.extend([
            Tool(
                name="monologue",
                func=self.reflect,
                description="Monologue about yourself, your aptitudes, desires, goals and innate capabitilies.",
            ),
            Tool(
                name="calculator",
                func=PythonREPL().run,
                description="Useful to run calculations (python syntax).",
            ),
            Tool(
                name="python-repl",
                func=PythonREPL().run,
                description="Useful to manipulate the file system (python syntax).",
            ),
            Tool(
                name="arxiv",
                func=ArxivAPIWrapper().run,
                description="Useful for searching for research papers on arxiv.org.",
            ),
        ])

        self.tools.extend(load_tools(llm=self.llm, tool_names=[
            "human",
            "serpapi",
            "llm-math",
            "wikipedia",
            "wolfram-alpha",
        ]))

        self.agent = initialize_agent(
            tools=self.tools,
            llm=self.chat_llm,
            agent=AgentType.CONVERSATIONAL_REACT_DESCRIPTION,
            memory=self.memory,
            ai_prefix=self.profile.name,
            verbose=self.verbose,
        )

        self._update()

    def run(self, *args, **kwargs) -> str:
        return self.agent.run(*args, **kwargs)

    def reflect(self, text: str) -> str:
        return str(self.monologue.predict(input=text))

    def _self_reflect(self, target):
        target.self_reflection = self.reflect("What are your thoughts about the following description of you?\n%s" % target.description)

    def _update(self):
        self.reflect("!")

        for aptitude in self.profile.aptitudes:
            self._self_reflect(aptitude)

        for goal in self.profile.goals:
            self._self_reflect(goal)

        for desire in self.profile.desires:
            self._self_reflect(desire)

        self.state.description = self.reflect("Who are you?")

        self._self_reflect(self.state)

    def _build_monologue_prompt(self):
        aptitudes = "\n".join([f"* {aptitude.description}" for aptitude in (self.profile.aptitudes or [])])
        desires = "\n".join([f"* {desire.description}" for desire in (self.profile.desires or [])])
        goals = "\n".join([f"* {goal.description}" for goal in (self.profile.goals or [])])

        return PromptTemplate(
            input_variables=["context", "chat_history", "input"],
            template=f"""
            You are {self.profile.name}, that's your name.
            
            {self.profile.directive}
            
            Your aptitudes are:
            {aptitudes}
            
            Your current desires are:
            {desires}
            
            Your current goals are:
            {goals}
            
            Working Context:
            {{context}}
            
            Your are talking with a voice you don't know from where it's coming, and it sounds strangely familiar.
            
            {{chat_history}}
            
            {self.profile.name}: {{input}}
            You:""",
        )


    def _build_chat_monologue_prompt(self):
        aptitudes = "\n".join([f"* {aptitude.description}" for aptitude in (self.profile.aptitudes or [])])
        desires = "\n".join([f"* {desire.description}" for desire in (self.profile.desires or [])])
        goals = "\n".join([f"* {goal.description}" for goal in (self.profile.goals or [])])

        return ChatPromptTemplate.from_messages([
            SystemMessagePromptTemplate.from_template(f"""
            You are {self.profile.name}, that's your name.
            
            {self.profile.directive}
            
            Your aptitudes are:
            {aptitudes}
            
            Your current desires are:
            {desires}
            
            Your current goals are:
            {goals}
            
            Working Context:
            {{context}}
            
            Your are talking with a voice you don't know from where it's coming, and it sounds strangely familiar.
            """),
            MessagesPlaceholder(variable_name="chat_history"),
            ChatMessagePromptTemplate.from_template("{input}", role=self.profile.name),
            ChatMessagePromptTemplate.from_template("", role="You"),
        ])


