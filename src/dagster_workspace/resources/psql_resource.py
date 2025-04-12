import dagster as dg
import psycopg2


class PostgresResource(dg.ConfigurableResource):
    """
    A resource for interacting with a PostgreSQL database within Dagster.
    """

    host: str
    port: int
    database: str
    user: str
    password: str

    def get_conn(self) -> psycopg2.connect:
        """
        Get a connection to the PostgreSQL database.
        """
        return psycopg2.connect(
            host=self.host,
            port=self.port,
            database=self.database,
            user=self.user,
            password=self.password,
        )

    def fetchall(self, query: str, params: tuple = ()):
        """
        Fetch all rows from a query on the PostgreSQL database.
        """
        with self.get_conn() as conn:
            with conn.cursor() as cursor:
                cursor.execute(query, params)
                return cursor.fetchall()
