import os

import pinecone
from langchain.vectorstores import Pinecone
from langchain.embeddings.openai import OpenAIEmbeddings

class Indexer:
    def __init__(
            self,
            index_name,
            namespace=None,
    ):
        self.embeddings = OpenAIEmbeddings()
        self.index_name = index_name
        self.namespace = namespace
        self.vs = None

    def create_index(self):
        pinecone.create_index(
            name=self.index_name,
            dimension=1536,  # dimensionality of dense model
            metric="dotproduct",  # sparse values supported only for dotproduct
            pod_type="s1",
            metadata_config={"indexed": []}  # see explaination above
        )

    def get_vector_store(self):
        if self.vs is None:
            self.vs = Pinecone.from_existing_index(
                embedding=self.embeddings,
                index_name=self.index_name,
                namespace=self.namespace,
            )

        return self.vs

    def as_retriever(self, **kwargs):
        return self.get_vector_store().as_retriever(**kwargs)

    def import_documents(self, documents):
        Pinecone.from_documents(
            documents=documents,
            embedding=self.embeddings,
            index_name=self.index_name,
            namespace=self.namespace,
        )

    def truncate(self):
        index = pinecone.Index(self.index_name)

        index.delete(
            delete_all = True,
            namespace = self.namespace,
        )


if os.environ["PINECONE_API_KEY"] is not None:
    pinecone.init(api_key=os.environ["PINECONE_API_KEY"], environment=os.environ["PINECONE_ENV"])
