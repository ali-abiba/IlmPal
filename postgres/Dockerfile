FROM postgres:15

# Environment variables for database configuration
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DB=postgres

# Create a directory for initialization scripts
RUN mkdir -p /docker-entrypoint-initdb.d

# Copy schema.sql to be executed on container initialization
COPY ./schema.sql /docker-entrypoint-initdb.d/01-schema.sql

# Expose the PostgreSQL port
EXPOSE 5432

# Optional: Custom PostgreSQL configuration
# COPY postgresql.conf /etc/postgresql/postgresql.conf
# CMD ["postgres", "-c", "config_file=/etc/postgresql/postgresql.conf"]

# Set the data directory
VOLUME ["/var/lib/postgresql/data"]

# Health check
HEALTHCHECK --interval=10s --timeout=5s --start-period=10s --retries=3 \
  CMD pg_isready -U postgres || exit 1
