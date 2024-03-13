FROM oven/bun

# Make the container discoverable by Styx and expose the expected port
LABEL me.saitho.styx.service=1

WORKDIR /app
COPY . ./
RUN bun install
CMD bun start