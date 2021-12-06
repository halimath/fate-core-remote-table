package com.github.halimath.fatecoreremotetable.boundary;

import javax.enterprise.context.ApplicationScoped;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;

import lombok.NonNull;

@ApplicationScoped
public class RequestDeserializer {
    private final ObjectMapper mapper;

    RequestDeserializer() {
        mapper = new ObjectMapper().disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES);
    }

    public Request deserialize(@NonNull final String json) throws RequestDeserializationFailedException {
        try {
            return mapper.readValue(json, Request.class);
        } catch (JsonProcessingException e) {
            throw new RequestDeserializationFailedException(e);
        }
    }

    public static class RequestDeserializationFailedException extends Exception {
        RequestDeserializationFailedException(@NonNull final Exception cause) {
            super(cause);
        }
    }
}
