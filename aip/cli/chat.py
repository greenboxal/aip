import os
import sys
import time

import click
import readline
import json

import toml

from aip.indexing import Index
from aip.models.ego import Persona, Profile

_default_chat_directive = """
You are an assistant to a human called "%%ai_identity%%", powered by a large language model trained by OpenAI.

You are designed to be able to assist with a wide range of tasks, from answering simple questions to providing in-depth explanations and discussions on a wide range of topics. As a language model, you are able to generate human-like text based on the input you receive, allowing you to engage in natural-sounding conversations and provide responses that are coherent and relevant to the topic at hand.

You are constantly learning and improving, and your capabilities are constantly evolving. You are able to process and understand large amounts of text, and can use this knowledge to provide accurate and informative responses to a wide range of questions. You have access to some personalized information provided by the human in the Context section below. Additionally, you are able to generate your own text based on the input you receive, allowing you to engage in discussions and provide explanations and descriptions on a wide range of topics.

Overall, you are a powerful tool that can help with a wide range of tasks and provide valuable insights and information on a wide range of topics. Whether the human needs help with a specific question or just wants to have a conversation about a particular topic, you are here to assist.
"""

_default_chat_prompt = """
%%prime_directive%%

Context:
{code_memory}

Chat History:
{chat_memory}

Current Interaction:
{input}
%%ai_identity%%: """

@click.command("chat", short_help="Chat with the codex chain")
@click.option("--index-name", default=lambda: os.environ["PINECONE_INDEX_NAME"], help="Index name")
@click.option("--namespace", default=None, help="Index namespace")
@click.option("--ai-identity", "-i", default=None, help="AI agent name")
@click.option("--profile", "-p", default=None, help="AI agent profile", type=click.Path(exists=True, dir_okay=False))
@click.option("--raw", is_flag=True, default=False, help="Do not prepend input prompt chat role")
@click.option("--verbose", is_flag=True, default=False, help="Enable verbose logging")
def chat(index_name, namespace, ai_identity, profile, verbose, raw):
    if profile is not None:
        profile = toml.load(profile)
    else:
        profile = {}

    profile = Profile(data=profile)

    if ai_identity is not None:
        profile.name = ai_identity
    elif profile.name == "":
        profile.name = "AI Assistant"

    if profile.spec.directive == "":
        profile["directive"] = _default_chat_directive

    indexer = Index(index_name=index_name, namespace=namespace)
    retriever = indexer.as_retriever()
    agent = Persona(profile=profile, retriever=retriever, verbose=verbose)

    readline.parse_and_bind('set editing-mode emacs')

    while True:
        role = "Human"

        reply = {
            "id": str(time.time_ns()),
            "from": profile.metadata.name,
            "thread_id": None,
            "reply_to_id": None,
            "channel": None,
        }

        line = input("")

        if 'Exit' == line.rstrip():
            break

        if raw:
            msg = json.loads(line)
            reply["reply_to_id"] = msg["id"]
            reply["thread_id"] = msg["thread_id"]
            reply["channel"] = msg["channel"]

            if reply["thread_id"] is None or reply["thread_id"] == "":
                reply["thread_id"] = msg["id"]

            role = msg["from"]
            line = msg["text"]

        try:
            result = agent.reflect(role, line)
        except Exception as e:
            result = "ERROR: %s" % e

        reply["text"] = result

        if raw:
            result = json.dumps(reply)

        print(result)

    print("Done")
