import asyncio
import os
import sys
import time

import aiofiles
import click
import json
import toml
from jsonrpc import JSONRPCResponseManager, Dispatcher

from aip.indexing import Index
from aip.models.ego import Persona, Profile

@click.command("ipc", short_help="Run IPC server")
@click.option("--index-name", default=lambda: os.environ["PINECONE_INDEX_NAME"], help="Index name")
@click.option("--namespace", default=None, help="Index namespace")
@click.option("--ai-identity", "-i", default=None, help="AI agent name")
@click.option("--profile", "-p", default=None, help="AI agent profile", type=click.Path(exists=True, dir_okay=False))
@click.option("--verbose", is_flag=True, default=False, help="Enable verbose logging")
def ipc(index_name, namespace, ai_identity, profile, verbose):
    if profile is not None:
        with open(profile, "r") as f:
            profile = json.load(f)
    else:
        profile = {}

    profile = Profile(data=profile)

    if ai_identity is not None:
        profile.name = ai_identity
    elif profile.name == "":
        profile.name = "AI Assistant"

    indexer = Index(index_name=index_name, namespace=namespace)
    retriever = indexer.as_retriever()
    agent = Persona(profile=profile, retriever=retriever, verbose=verbose)

    def route_message(msg):
        reply = {
            "id": str(time.time_ns()),
            "from": profile.name,
            "thread_id": None,
            "reply_to_id": None,
            "channel": None,
        }

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

        return reply

    ipc_base_fd = int(os.environ["AIP_IPC_BASE_FD"])
    ipc_in_fd = ipc_base_fd
    ipc_out_fd = ipc_base_fd + 1

    with os.fdopen(ipc_in_fd, "rb") as ipc_in:
        with os.fdopen(ipc_out_fd, "wb") as ipc_out:
            def handle(line=None):
                if line is None:
                    return

                request = line.strip()

                if request == "":
                    return

                request = json.loads(request)
                result = route_message(request)

                ipc_out.write(bytes(json.dumps(result) + "\n", 'utf-8'))
                ipc_out.flush()

            for line in ipc_in:
                handle(line)

        #loop = asyncio.get_event_loop()
        #loop.add_reader(sys.stdin, handle)
        #loop.run_forever()
