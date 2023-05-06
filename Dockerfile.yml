FROM node:alpine

RUN apk update && apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont

ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD true
ENV PUPPETEER_EXECUTABLE_PATH /usr/bin/chromium-browser

RUN npm install -g npm@latest
RUN npm install -g wappalyzer@latest

ENTRYPOINT ["/usr/local/bin/wappalyzer"]
CMD ["-P"]