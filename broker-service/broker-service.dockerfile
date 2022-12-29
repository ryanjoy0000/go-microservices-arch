# # set the base go image and set a build stage
# FROM golang:1.19.4-alpine as builder

# # create folder app in root
# RUN mkdir /app

# # copy files from current service folder to app folder 
# COPY . /app

# # set the working directory
# WORKDIR /app

# # set env variable, build the go code and set output for executable
# RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# # make sure it is executable by setting permissions
# RUN chmod +x /app/brokerApp

# # --2nd stage--

# # build a new tiny docker image
# FROM alpine

# # create folder in root
# RUN mkdir /app

# # copy executable from previous image to new image
# COPY --from=builder /app/brokerApp /app

# # provide defaults for an executing container
# CMD ["/app/brokerApp"]

#----------------------------------------------------------------------------
#----------------------------------------------------------------------------

# With makefile

# build a new tiny docker image
FROM alpine

# create folder in root
RUN mkdir /app

# copy executable from previous image to new image
COPY brokerApp /app

# provide defaults for an executing container
CMD ["/app/brokerApp"]