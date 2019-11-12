#!/usr/bin/env python
# coding: utf-8
from __future__ import unicode_literals, print_function

import datetime
import os
import re
from multiprocessing.pool import ThreadPool
from subprocess import PIPE, Popen

from PIL import Image
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

    def create_thumb(self):
        conf = self.__get_conf()
        image_dir = conf["static"]["MEDIA_ROOT"]
        width = 160
        height = 160
        regex = re.compile('(?!.*_thumb)(.*)\.(jpe?g|png|gif)$')
        regex2 = re.compile('^.*_thumb\.jpg$')
        images = []
        for root, dirs, files in os.walk(image_dir):
            files2 = list(filter(lambda x: regex2.match(x), files))
            for file in filter(lambda x: regex.match(x), files):
                name = regex.findall(file)[0][0]
                if not list(filter(lambda x: name in x, files2)):
                    images.append((root, file, name,))

        pool = ThreadPool(processes=10)
        def resize(val):
            root, file, name = val
            im1 = Image.open(os.path.join(root, file)).convert('RGB')
            im2 = im1.resize((width, height), Image.NEAREST)
            path = os.path.join(root, name + "_thumb" + ".jpg")
            im2.save(path)
            print("Created thumb {}".format(path))
        print("Count: {}".format(len(images)))
        pool.map(resize, images)

    def create_admins(self):
        conf = self.__get_conf()
        admins = conf["admins"]
        for admin in admins:
            name = admin["NAME"]
            first_name = name
            last_name = None
            try:
                first_name, last_name = name.split(" ")[:2]
            except ValueError:
                pass

            email = admin["EMAIL"]
            print("Input login from {} (enter 'c' to skip):".format(name))
            login = input()
            if login == 'c':
                continue
            if not login:
                print("Login cannot be empty:".format(name))
                continue

            password = getpass.getpass("Input password:").encode("ascii")
            confirm_password = getpass.getpass("Confirm password:").encode("ascii")

            if password == confirm_password:
                db = self.__connect_default_db()
                cursor = db.cursor(cursor_factory=DictCursor)
                token = binascii.hexlify(os.urandom(16)).decode('ascii')
                password = bcrypt.hashpw(password, bcrypt.gensalt(4)).decode('ascii')
                cursor.execute(
                    "SELECT id FROM users WHERE  user_name = %s OR email = %s",
                    (login, email)
                )
                try:
                    user_id = cursor.fetchone()[0]
                except TypeError:
                    user_id = 0

                if user_id:
                    cursor.execute("""
                        DELETE FROM tokens WHERE id = (SELECT token_id FROM users WHERE id = %s)
                    """, (user_id,))

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
                if not user_id:
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
                    print("Admin {} created!".format(name))
                else:
                    cursor.execute(
                        """
                        UPDATE users 
                        SET user_name = %s, 
                            first_name = %s,
                            last_name = %s,
                            password = %s,
                            email = %s,
                            token_id = %s,
                            updated_at = now(),
                            deleted_at = NULL
                        WHERE id = %s
                        """,
                        (login, first_name, last_name, password, email, token_id, user_id)
                    )
                    print("Admin {} updated!".format(name))
                db.commit()
            else:
                print ("Passwords do not match!")

    def backup_db(self):
        conf = self.__get_conf()
        default_db_conf = conf["databases"]["default"]
        command = "PGPASSWORD='{password}' pg_dump -h localhost -p 5432 -U {user} -F c -b -v -f {backup_dir}/{backup_name} {database}"
        command = command.format(
            password=default_db_conf["PASSWORD"],
            host=default_db_conf["HOST"],
            user=default_db_conf["USER"],
            database=default_db_conf["DATABASE"],
            backup_dir=conf["system"]["BACKUP_PATH"],
            backup_name="{}-{}.backup".format(
                datetime.datetime.now().strftime("%Y%m%d"),
                default_db_conf["DATABASE"]
            )
        )
        p = Popen(command, shell=True, stdin=PIPE)
        print(p.stdout)

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
