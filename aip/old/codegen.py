import pathlib
import re
import shutil

from langchain.chains import ConversationChain
from langchain.memory import (
    ConversationSummaryBufferMemory,
    ConversationEntityMemory,
    ConversationBufferWindowMemory,
    CombinedMemory,
)
from langchain.chat_models import ChatOpenAI
from langchain.prompts import PromptTemplate

BEGIN_FILE_MARKER = re.compile(r'OUTPUT FILE +([a-zA-Z0-9_./ -]+):')
END_FILE_MARKER = re.compile(r'LLM EOF')
EOM_MARKER = re.compile(r'END OF MESSAGE')

_objective = """
You're an expert in Geometric Algebra, Linear Algebra, and Computer Science, in special in Go.

You are tasked to write a complete Go package called "gealg" for Geometric Algebra satisfying the following requirements:

* Support Pseudoscalars, Bivectors, Trivectors and Pseudovectors on the Clifford algebra Cl(3,0,1)
* Support for the Wedge product using a Cayley table
* Define constants for the standard basis vectors
* Optimized for performance

Split each declared type and its functions in different files. List me the output directory structure and stand by to output each file individual content upon request.
Do not leave any TODOs or any missing implementations.
Annotate each file with a comment on its very first line in the following format: "// OUTPUT FILE <filename>:"
Show the entire content of each file upon request.
Annotate each file with a comment on its very last line in the following format: "// LLM EOF"
Say "END OF MESSAGE" after all files.
"""

_chat_prompt = """
You are an assistant to a human, powered by a large language model trained by OpenAI.

You are designed to be able to assist with a wide range of tasks, from answering simple questions to providing in-depth explanations and discussions on a wide range of topics. As a language model, you are able to generate human-like text based on the input you receive, allowing you to engage in natural-sounding conversations and provide responses that are coherent and relevant to the topic at hand.

You are constantly learning and improving, and your capabilities are constantly evolving. You are able to process and understand large amounts of text, and can use this knowledge to provide accurate and informative responses to a wide range of questions. You have access to some personalized information provided by the human in the Context section below. Additionally, you are able to generate your own text based on the input you receive, allowing you to engage in discussions and provide explanations and descriptions on a wide range of topics.

Overall, you are a powerful tool that can help with a wide range of tasks and provide valuable insights and information on a wide range of topics. Whether the human needs help with a specific question or just wants to have a conversation about a particular topic, you are here to assist.

Long Memory:
{long_term_memory}

Context:
{entities}

Short Memory:
{short_term_memory}

Chat History:
{chat_memory}
Human: {input}
AI Assistant: """

chat_prompt = PromptTemplate(
    input_variables=["entities", "chat_memory", "input", "long_term_memory", "short_term_memory"],
    template=_chat_prompt,
)

initial_prompt = _objective + """

What is the content of each file? Show me the entire content of each file upon request.
"""

continue_prompt = """Continue."""


def main():
    llm_feeling = ChatOpenAI(model_name="gpt-3.5-turbo", temperature=0.30, max_tokens=128)
    llm_memory = ChatOpenAI(model_name="gpt-3.5-turbo", temperature=0.70, max_tokens=512)
    llm_codex = ChatOpenAI(model_name="gpt-3.5-turbo", temperature=0.70, max_tokens=256)
    llm_reason = ChatOpenAI(model_name="gpt-4-32k", temperature=0.70, max_tokens=16000)

    short_term_memory = ConversationSummaryBufferMemory(llm=llm_memory, max_token_limit=200, memory_key="short_term_memory", input_key="input", ai_prefix = "AI Assistant")
    long_term_memory = ConversationSummaryBufferMemory(llm=llm_feeling, max_token_limit=500, memory_key="long_term_memory", input_key="input", ai_prefix = "AI Assistant")
    chat_memory = ConversationBufferWindowMemory(memory_key="chat_memory", input_key="input", ai_prefix="AI Assistant", k=5)
    entity_memory = ConversationEntityMemory(llm=llm_memory, ai_prefix="AI Assistant", chat_history_key="chat_memory", input_key="input")
    #graph_memory = ConversationKGMemory(llm=llm_feeling, ai_prefix="AI Assistant", memory_key="graph_memory", input_key="input")

    memory = CombinedMemory(memories=[
        short_term_memory,
        long_term_memory,
        chat_memory,
        entity_memory,
        #graph_memory,
    ])

    codex_chain = ConversationChain(
        llm=llm_codex,
        verbose=True,
        memory=memory,
        prompt=chat_prompt,
    )

    counter = 0
    parser = FileSplitter("./output", clean=True, immediate=True)

    while True:
        if counter == 0 :
            result = codex_chain.predict(input=initial_prompt)
        else:
            result = codex_chain.predict(input=continue_prompt)

        for line in result.splitlines():
            parser.parse(line)

        counter += 1

        if re.search(EOM_MARKER, result) is not None and counter > 5:
            break

    parser.emit()

    print("Completed in %d iterations" % counter)


class FileSplitter():
    def __init__(self, output_path, immediate=False, clean=False):
        self.immediate = immediate
        self.current_file = ""
        self.output_path = pathlib.Path(output_path)
        self.files = {}

        if clean:
            shutil.rmtree(str(self.output_path), ignore_errors=True)

        self.output_path.mkdir(parents=True, exist_ok=True)

    def parse(self, line):
        print(line)

        m = re.search(BEGIN_FILE_MARKER, line)

        if m is not None:
            self.begin_file(m.group(1))
            return

        m = re.search(END_FILE_MARKER, line)

        if m is not None:
            self.end_file()
            return

        self.append_line(line)

    def begin_file(self, name):
        if self.current_file != name:
            self.end_file()

        self.current_file = name

    def end_file(self):
        if self.current_file == "":
            return

        if self.immediate:
            self.emit_file(self.current_file)

        self.current_file = ""

    def append_line(self, line):
        if self.current_file != "":
            if self.current_file not in self.files:
                self.files[self.current_file] = ""

            self.files[self.current_file] += line + "\n"


    def emit_file(self, name):
        if not name in self.files:
            return

        contents = self.files[name]
        path = self.output_path.joinpath(name)

        path.parent.mkdir(parents=True, exist_ok=True)

        with path.open(mode="w") as f:
            f.write(contents)

    def emit(self):
        self.end_file()

        if not self.immediate:
            for file in self.files.keys():
                self.emit_file(file)


main()
