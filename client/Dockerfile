FROM oven/bun

# Make the container discoverable by Styx and expose the expected port
LABEL me.saitho.styx.service=1
# default port is 8844. Set different port in label to let Styx know
# make sure your application servers the default port 8844 or the port set via label
LABEL me.saitho.styx.port=8843

# environment for app, see `process.env.PORT` in index.ts
ENV PORT=8843

WORKDIR /app
COPY . ./
RUN bun install
CMD bun start