#!/usr/bin/env python
# -*- coding: utf-8 -*-

import re

from bson import ObjectId
from pymongo import MongoClient
from pymongo.errors import ServerSelectionTimeoutError
import sys
import pprint
import string
import random

def get_random_string():
    return ''.join(random.choices(string.ascii_lowercase, k=10))

mongo_server = "localhost:27017"

client = MongoClient(mongo_server, serverSelectionTimeoutMS=200, unicode_decode_error_handler='ignore')
db = client.pcap
pcap_coll = db.pcap
files_coll = db.filesImported

random_string = get_random_string()
inp = input(f"Are you sure? Enter {random_string} to empty the pcap database: ")
if inp == random_string:
    pcap_coll.drop()
    files_coll.drop()
    print("pcap database dropped")
