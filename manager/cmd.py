#!/usr/bin/env python
# coding: utf-8
from __future__ import unicode_literals, print_function
import os

from fire import Fire
import yaml
import getpass
import psycopg2
import bcrypt
import binascii
from builtins import input
from psycopg2.extras import DictCursor


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
                (created_at, updated_at, first_name, last_name, user_name, password, email, token_id, is_super_user, is_staff)
                VALUES
                (now(), now(), %s, %s, %s, %s, %s, %s, TRUE, TRUE)
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

    def init_data(self):
        pass


if __name__ == '__main__':
    Fire(Commands)
