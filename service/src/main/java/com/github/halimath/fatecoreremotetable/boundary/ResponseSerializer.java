package com.github.halimath.fatecoreremotetable.boundary;

import javax.enterprise.context.ApplicationScoped;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.halimath.fatecoreremotetable.boundary.dto.Response;

import lombok.NonNull;

@ApplicationScoped
class ResponseSerializer {
    private final ObjectMapper mapper;

    ResponseSerializer() {
        mapper = new ObjectMapper();
    }

    String serialize (@NonNull final Response response) {
        try {
            return mapper.writeValueAsString(response);
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        }
    }    
}
