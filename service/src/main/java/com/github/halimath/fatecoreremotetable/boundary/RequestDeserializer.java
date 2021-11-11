package com.github.halimath.fatecoreremotetable.boundary;

import javax.enterprise.context.ApplicationScoped;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.halimath.fatecoreremotetable.boundary.dto.Request;

import lombok.NonNull;

@ApplicationScoped
class RequestDeserializer {
    private final ObjectMapper mapper;

    RequestDeserializer() {
        mapper = new ObjectMapper().disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES);
    }

    Request deserialize(@NonNull final String json) throws RequestDeserializationFailedException {
        try {
            return mapper.readValue(json, Request.class);
        } catch (JsonProcessingException e) {
            throw new RequestDeserializationFailedException(e);
        }
    }

    static class RequestDeserializationFailedException extends Exception {
        RequestDeserializationFailedException(@NonNull final Exception cause) {
            super(cause);
        }
    }
}
