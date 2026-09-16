# agent-workbench.buildCommand: ./bazel.cmd build //...
FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends \
    git curl ca-certificates ripgrep patch && \
    rm -rf /var/lib/apt/lists/*
RUN ARCH=$(dpkg --print-architecture) && \
    JBR_ARCH=$([ "$ARCH" = "amd64" ] && echo "x64" || echo "aarch64") && \
    curl -Lo /tmp/jbr.tar.gz \
    "https://cache-redirector.jetbrains.com/intellij-jbr/jbr_jcef-25.0.2-linux-${JBR_ARCH}-b329.111.tar.gz" && \
    mkdir -p /opt/jbr && tar -xzf /tmp/jbr.tar.gz -C /opt/jbr --strip-components=1 && \
    rm /tmp/jbr.tar.gz
ENV JAVA_HOME=/opt/jbr
ENV PATH=$JAVA_HOME/bin:$PATH
RUN ARCH=$(dpkg --print-architecture) && \
    curl -Lo /usr/local/bin/bazel \
    https://github.com/bazelbuild/bazelisk/releases/latest/download/bazelisk-linux-${ARCH} && \
    chmod +x /usr/local/bin/bazel
RUN echo 'common --disk_cache=/bazel-cache' > /etc/bazel.bazelrc
RUN echo 'startup --host_jvm_args=-Xmx8g' > /root/.bazelrc
RUN useradd -m -s /bin/bash agent
WORKDIR /workspace
