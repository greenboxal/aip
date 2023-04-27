import os
import sys

import click
import readline

from aip.indexing import Indexer
from aip.llm import Agent


@click.command("chat")
@click.option("--index-name", default=lambda: os.environ["PINECONE_INDEX_NAME"])
@click.option("--namespace", default=None)
@click.option("--model", default="gpt-3.5-turbo")
def chat(index_name, namespace, model):
    print("Using model {}".format(model))
    indexer = Indexer(index_name=index_name, namespace=namespace)
    retriever = indexer.as_retriever(search_type="similarity", search_kwargs={"k":5})
    agent = Agent(retriever, model)

    readline.parse_and_bind('set editing-mode emacs')

    while True:
        line = input(">>> ")

        if 'Exit' == line.rstrip():
            break

        try:
            result = agent.codex_chain.predict(input=line)
            print(result)
        except Exception as err:
            print("### EXCEPTION ###")
            print(err)

    print("Done")
