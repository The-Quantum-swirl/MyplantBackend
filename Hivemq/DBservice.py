import psycopg2
from datetime import datetime, timezone

class DBService:

    def __init__(self):
        print("starting db service")
    
    def start_connection(self):
        self.conn = psycopg2.connect(
            database="core-service",
            user='postgres',
            password='quicuxeo',
            host='localhost',
            port='5432'
        )
    
    def save(self):
        # Commit your changes in the database
        self.conn.commit()
        # Closing the connection
        self.conn.close()


    def insert_data(self, email, deviceId):
        self.start_connection()
        deviceType = "none"
        if "wp-" in deviceId:
            deviceType = "wp"

        dt = datetime.now(timezone.utc)

        data = (email, deviceId, deviceType, "", "", str(dt), "", False, "")

        print("inserting "+str(data))
        # creating a cursor object
        cursor = self.conn.cursor()
        sql = "INSERT INTO public.user VALUES(%s, %s, %s, %s, %s, %s, %s, %s, %s)"
        cursor.execute(sql, data)

        self.save()
    
    def fetchDeviceId(self, email):
        self.start_connection()

        # creating a cursor object
        cursor = self.conn.cursor()
        sql = "SELECT * FROM public.user WHERE email = '" + email +"'"
        print("executing :" + sql)
        cursor.execute(sql)
        obj = cursor.fetchone();

        # Closing the connection
        self.conn.close()
        return obj
