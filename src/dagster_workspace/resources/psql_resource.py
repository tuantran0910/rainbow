import dagster as dg
import psycopg2


class PostgresResource(dg.ConfigurableResource):
    """
    A resource for interacting with a PostgreSQL database within Dagster.

    Args:
        host (str): The host of the PostgreSQL database.
        port (int): The port of the PostgreSQL database.
        database (str): The name of the PostgreSQL database.
        user (str): The user of the PostgreSQL database.
        password (str): The password of the PostgreSQL database.
    """

    host: str
    port: int
    database: str
    user: str
    password: str

    def get_conn(self) -> psycopg2.connect:
        """
        Get a connection to the PostgreSQL database.

        Returns:
            psycopg2.connect: A connection to the PostgreSQL database.
        """
        return psycopg2.connect(
            host=self.host,
            port=self.port,
            database=self.database,
            user=self.user,
            password=self.password,
        )

    def fetchall(self, query: str, params: tuple = ()) -> list[tuple]:
        """
        Fetch all rows from a query on the PostgreSQL database.

        Args:
            query (str): The query to execute.
            params (tuple): The parameters to pass to the query.

        Returns:
            list[tuple]: A list of tuples containing the rows from the query.
        """
        with self.get_conn() as conn:
            with conn.cursor() as cursor:
                cursor.execute(query, params)
                return cursor.fetchall()
