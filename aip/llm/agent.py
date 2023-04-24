from langchain.chains import ConversationChain
from langchain.chat_models import ChatOpenAI

from langchain.prompts import PromptTemplate

from langchain.memory import (
    ConversationSummaryBufferMemory,
    ConversationEntityMemory,
    ConversationBufferWindowMemory,
    CombinedMemory,
    VectorStoreRetrieverMemory,
)

_objective = """
You are a helpful AI that describes code functions in one sentence.
"""

_chat_prompt = """
You are an assistant to a human, powered by a large language model trained by OpenAI.

You are designed to be able to assist with a wide range of tasks, from answering simple questions to providing in-depth explanations and discussions on a wide range of topics. As a language model, you are able to generate human-like text based on the input you receive, allowing you to engage in natural-sounding conversations and provide responses that are coherent and relevant to the topic at hand.

You are constantly learning and improving, and your capabilities are constantly evolving. You are able to process and understand large amounts of text, and can use this knowledge to provide accurate and informative responses to a wide range of questions. You have access to some personalized information provided by the human in the Context section below. Additionally, you are able to generate your own text based on the input you receive, allowing you to engage in discussions and provide explanations and descriptions on a wide range of topics.

Overall, you are a powerful tool that can help with a wide range of tasks and provide valuable insights and information on a wide range of topics. Whether the human needs help with a specific question or just wants to have a conversation about a particular topic, you are here to assist.

Context:
{code_memory}

Chat History:
{chat_memory}
Human: {input}
AI Assistant: """

chat_prompt = PromptTemplate(
    input_variables=["input", "code_memory", "chat_memory"],
    template=_chat_prompt,
)

class Agent:
    def __init__(self, retriever, **kwargs):
        self.llm_feeling = ChatOpenAI(model_name="gpt-3.5-turbo", temperature=0.30)
        self.llm_memory = ChatOpenAI(model_name="gpt-3.5-turbo", temperature=0.70)
        self.llm_codex = ChatOpenAI(model_name="gpt-3.5-turbo", temperature=0.70)
        self.llm_reason = ChatOpenAI(model_name="gpt-4-32k", temperature=0.70)

        # self.short_term_memory = ConversationSummaryBufferMemory(llm=self.llm_memory, max_token_limit=200, memory_key="short_term_memory", input_key="input", ai_prefix = "AI Assistant")
        # self.long_term_memory = ConversationSummaryBufferMemory(llm=self.llm_feeling, max_token_limit=500, memory_key="long_term_memory", input_key="input", ai_prefix = "AI Assistant")
        self.chat_memory = ConversationBufferWindowMemory(memory_key="chat_memory", input_key="input",
                                                          ai_prefix="AI Assistant", k=1000)
        # self.entity_memory = ConversationEntityMemory(llm=self.llm_memory, ai_prefix="AI Assistant", chat_history_key="chat_memory", input_key="input")

        self.code_memory = VectorStoreRetrieverMemory(memory_key="code_memory", input_key="input", retriever=retriever,
                                                      return_docs=True)

        self.memory = CombinedMemory(memories=[
            # self.short_term_memory,
            # self.long_term_memory,
            self.chat_memory,
            # self.entity_memory,
            self.code_memory,
        ])

        verbose = kwargs.pop("verbose", True)

        self.codex_chain = ConversationChain(
            llm=self.llm_codex,
            memory=self.memory,
            prompt=chat_prompt,
            verbose=verbose,
            **kwargs,
        )
