import os

import click

from aip.indexing import Indexer

from langchain.document_loaders import (
    DirectoryLoader,
    TextLoader,
)

from langchain.text_splitter import (
    TokenTextSplitter,
)

@click.command("create")
@click.argument("name", default=lambda: os.environ["PINECONE_INDEX_NAME"])
def index_create(name):
    indexer = Indexer(index_name=name)

    print("Creating index %s" % name)
    indexer.create_index()

    print("Index %s created" % name)

@click.command("truncate")
@click.argument("name", default=lambda: os.environ["PINECONE_INDEX_NAME"])
@click.option("--namespace", default=None)
def index_truncate(name, namespace):
    indexer = Indexer(index_name=name)

    print("Truncating index %s (ns = %s)" % (name, namespace))
    indexer.truncate()

    print("Index %s truncated (ns = %s)" % (name, namespace))

@click.command("add")
@click.argument("file")
@click.option("--index-name", "-i", default=lambda: os.environ["PINECONE_INDEX_NAME"])
@click.option("--namespace", "-n", default=None)
@click.option("--filter", "-f", default="**/*")
@click.option("--relative-to", "-C", default="")
@click.option("--prefix", "-p", default="")
@click.option("--dry-run", is_flag=True, default=False)
@click.option("--chunk-size", default=2048)
@click.option("--chunk-overlap", default=0)
def index_add(index_name, relative_to, prefix, dry_run, filter, chunk_size, chunk_overlap, namespace, file):
    if index_name == "":
        index_name = os.environ["PINECONE_INDEX_NAME"]

    indexer = Indexer(index_name=index_name, namespace=namespace)

    if not os.path.exists(file):
        raise ValueError("File does not exist: %s" % file)

    if os.path.isfile(file):
        loader = TextLoader(file, encoding='utf-8')
    else:
        loader = DirectoryLoader(file, glob=filter, loader_cls=TextLoader, loader_kwargs={"encoding": "utf-8"})

    docs = loader.load()

    for doc in docs:
        if "source" in doc.metadata:
            source = doc.metadata["source"]
            source = os.path.relpath(source, relative_to)
            source = os.path.join(prefix, source)
            doc.metadata["source"] = source

        print(doc.metadata)

    text_splitter = TokenTextSplitter(chunk_size=chunk_size, chunk_overlap=chunk_overlap)
    texts = text_splitter.split_documents(docs)

    print("Importing %d documents (%d chunks)" % (len(docs), len(texts)))

    if dry_run:
        print("Dry run, not importing documents")
    else:
        indexer.import_documents(texts)
        print("Done")

@click.group("index")
def index():
    pass

index.add_command(index_create)
index.add_command(index_truncate)
index.add_command(index_add)
