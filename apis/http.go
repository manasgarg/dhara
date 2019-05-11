package apis

import (
	"fmt"
	"strconv"

	"github.com/buaazp/fasthttprouter"
	"github.com/manasgarg/dhara/stream"
	"github.com/manasgarg/dhara/utils"
	"github.com/valyala/fasthttp"
)

func postMessage(ctx *fasthttp.RequestCtx) {
	streamId := ctx.UserValue("streamId").(string)
	messageBody := ctx.PostBody()
	if messageBody == nil {
		ctx.Error("Must send a message body", fasthttp.StatusBadRequest)
		return
	}

	s, err := stream.GetMemStream(streamId)
	if err != nil {
		utils.SLogger.Error("Error in getting stream by id", streamId)
		ctx.Error("", fasthttp.StatusInternalServerError)
	} else {
		messageId := s.AddMessage(messageBody)
		fmt.Fprintf(ctx, "message_id=%d-%d", messageId.PartOne, messageId.PartTwo)
		utils.SLogger.Debugw("Added message", "stream", streamId, "message", messageId)
	}
}

func addMessage(ctx *fasthttp.RequestCtx) {
	arg := ctx.QueryArgs().Peek("stream_id")
	if arg == nil {
		ctx.Error("stream_id is a required parameter.", fasthttp.StatusBadRequest)
		return
	}

	streamId := string(arg)

	arg = ctx.QueryArgs().Peek("message")
	if arg == nil {
		ctx.Error("message is a required parameter.", fasthttp.StatusBadRequest)
		return
	}

	messageBody := make([]byte, len(arg))
	copy(messageBody, arg)

	if s, err := stream.GetMemStream(streamId); err != nil {
		utils.SLogger.Error("Error in getting stream by id", streamId)
		ctx.Error("", fasthttp.StatusInternalServerError)
	} else {
		messageId := s.AddMessage(messageBody)
		fmt.Fprintf(ctx, "message_id=%d-%d", messageId.PartOne, messageId.PartTwo)
		utils.SLogger.Debugw("Added message", "stream", streamId, "message", messageId)
	}
}

func getMessages(ctx *fasthttp.RequestCtx) {
	var err error

	streamId := ctx.UserValue("streamId").(string)
	startMessageId := stream.MessageId{PartOne: 0, PartTwo: 0}

	arg := ctx.QueryArgs().Peek("start_message_id")
	if arg != nil {
		_, err = fmt.Sscanf(string(arg), "%d-%d", &startMessageId.PartOne, &startMessageId.PartTwo)
		if err != nil {
			ctx.Error(fmt.Sprintf("start_message_id is in bad format. Supplied value [%s].",
				arg), fasthttp.StatusBadRequest)
			return
		}
	}

	count := 10
	arg = ctx.QueryArgs().Peek("count")
	if arg != nil {
		count, err = strconv.Atoi(string(arg))
		if err != nil {
			ctx.Error(fmt.Sprintf("count is in bad format. Supplied value [%s].",
				arg), fasthttp.StatusBadRequest)
			return
		}

		if count > 100 {
			count = 100
		}
	}

	delim := []byte{}
	arg = ctx.QueryArgs().Peek("delim")
	if arg != nil {
		delim = arg
	}

	s, err := stream.GetMemStream(streamId)
	if err != nil {
		utils.SLogger.Error("Error in getting stream by id", streamId)
		ctx.Error("", fasthttp.StatusInternalServerError)
		return
	}

	messages := s.GetMessages(&startMessageId, count)
	utils.SLogger.Debugw("Found messages", "stream", streamId, "start_message_id", startMessageId, "count", len(messages))
	for _, message := range messages {
		ctx.Write(message.Body)
		ctx.Write(delim)
	}
}

func StartHTTPServer(inf string, port int) error {
	router := fasthttprouter.New()
	router.POST("/stream/:streamId/messages", postMessage)
	router.GET("/add-messages", addMessage)
	router.GET("/stream/:streamId/messages", getMessages)

	return fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", inf, port), router.Handler)
}
