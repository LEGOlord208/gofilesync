FROM karalabe/xgo-latest
RUN apt update
RUN apt install -y libgtk-3-dev
RUN apt install -y libappindicator3-dev
