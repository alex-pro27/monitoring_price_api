#!/usr/bin/env python
# coding: utf-8
from __future__ import unicode_literals, print_function
import os
import re
from subprocess import PIPE, Popen
from fire import Fire
import yaml
import getpass
import psycopg2
import bcrypt
import binascii
from builtins import input
from psycopg2.extras import DictCursor

try:
    from os import scandir
except ImportError:
    from scandir import scandir  # use scandir PyPI module on Python < 3.5


class Commands(object):

    __conf = None
    __default_db = None

    def __del__(self):
        if self.__default_db:
            self.__default_db.close()

    def __get_conf(self):
        if not self.__conf:
            conf_path = os.environ["MONITORING_PRICE_CONF"]
            with open(conf_path) as stream:
                self.__conf = yaml.safe_load(stream)
        return self.__conf

    def __connect_default_db(self):
        if not self.__default_db:
            conf = self.__get_conf()
            self.__default_db = psycopg2.connect(
                "host='{HOST}' dbname='{DATABASE}' user='{USER}' password='{PASSWORD}' port={PORT}"
                .format(**conf["databases"]["default"])
            )
        return self.__default_db

    def create_admin(self):
        conf = self.__get_conf()
        admin = conf["admin"]["NAME"]
        first_name = admin
        last_name = None
        try:
            first_name, last_name = admin.split(" ")[:2]
        except ValueError:
            pass

        email = conf["admin"]["EMAIL"]
        print("Input login from {}:".format(admin))
        login = input()
        if not login:
            print("Login cannot be empty:".format(admin))
            return
        password = getpass.getpass("Input password:").encode("ascii")
        confirm_password = getpass.getpass("Confirm password:").encode("ascii")

        if password == confirm_password:
            db = self.__connect_default_db()
            cursor = db.cursor(cursor_factory=DictCursor)
            token = binascii.hexlify(os.urandom(16)).decode('ascii')
            password = bcrypt.hashpw(password, bcrypt.gensalt(4)).decode('ascii')
            cursor.execute(
                """
                INSERT INTO tokens 
                (created_at, updated_at, "key")
                VALUES 
                (now(), now(), %s)
                RETURNING id;
                """, (token,)
            )

            token_id = cursor.fetchone()[0]

            cursor.execute(
                """
                INSERT INTO users
                (created_at, updated_at, first_name, last_name, user_name, password, email, token_id, is_super_user)
                VALUES
                (now(), now(), %s, %s, %s, %s, %s, %s, TRUE)
                RETURNING id;
                """,
                (
                    first_name,
                    last_name,
                    login,
                    password,
                    email,
                    token_id,
                )
            )
            db.commit()
            print("Admin {} created!".format(admin))
        else:
            print ("Passwords do not match!")

    def dump(self):
        # TODO
        pass

    def init_data(self):
        command = "PGPASSWORD='{password}' pg_restore -c -h {host} -d {database} -U {user} < {filename}"
        data_path = os.path.join(os.path.realpath(os.path.dirname(__file__)), "data")
        files = scandir(data_path)
        conf = self.__get_conf()
        default_db_conf = conf["databases"]["default"]
        print("Can this command overwrite existing data, continue? (y/n)")
        ans = input()
        if ans not in ["Y", "y", "yes"]:
            print("canceled")
            return
        for f in files:
            if f.is_file() and re.match(".+\.(psql|du?mp|backup)$", f.name):
                n_command = command.format(
                    password=default_db_conf["PASSWORD"],
                    host=default_db_conf["HOST"],
                    user=default_db_conf["USER"],
                    database=default_db_conf["DATABASE"],
                    filename=f.path
                )
                print(n_command)
                p = Popen(n_command, shell=True, stdin=PIPE)
                print(p.stdout)

if __name__ == '__main__':
    Fire(Commands)
