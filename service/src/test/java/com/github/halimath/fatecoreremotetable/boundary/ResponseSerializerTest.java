package com.github.halimath.fatecoreremotetable.boundary;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.util.List;

import com.fasterxml.jackson.core.JsonProcessingException;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

public class ResponseSerializerTest {
    ResponseSerializer serializer;

    @BeforeEach
    void initSerializer() {
        serializer = new ResponseSerializer();
    }

    @Test
    void shouldSerializeErrorResponse() throws JsonProcessingException {
        final var got = serializer.serialize(
                new Response("1", "1", Response.Type.ERROR, null, new Response.Error(404, "Table not found")));
        assertEquals("""
                {"id":"1","self":"1","type":"error","error":{"code":404,"reason":"Table not found"}}""", 
                got);
    }

    @Test
    void shouldSerializeTableResponse() throws JsonProcessingException {
        final var got = serializer.serialize(new Response("1", "1", Response.Type.TABLE,
                new Response.Table("1", "test", "gm", List.of(new Response.Table.Player("1", "Cynere", 2, List.of())),
                        List.of(new Response.Table.Aspect("1", "Fog"))),
                null));

        assertEquals("""
            {"id":"1","self":"1","type":"table","table":{"id":"1","title":"test","gamemaster":"gm","players":[{"id":"1","name":"Cynere","fatePoints":2,"aspects":[]}],"aspects":[{"id":"1","name":"Fog"}]}}""",
            got);
    }
}
