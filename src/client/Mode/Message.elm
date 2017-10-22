module Mode.Message exposing (Message, view, query)

-- Core packages

import Http
import Html exposing (..)
import Html.Attributes exposing (id, style, class)
import Json.Decode as Decode exposing (Decoder, field, string, list)


-- Third-party packages

import GraphQl exposing (Operation, Query, Named)


-- MODEL


type alias Message =
    { posted : String
    , content : String
    }



-- VIEW


view : List Message -> Html msg
view messages =
    ul [ id "messageList", style [ ( "list-style", "none" ) ] ] <| List.map viewMessage <| List.reverse messages


viewMessage : Message -> Html msg
viewMessage { posted, content } =
    li [ class "adminMessage" ]
        [ span [ class "gray-text text-darken-4" ] [ text ("Posted " ++ posted) ]
        , p [] [ text content ]
        ]



-- QUERY


message : Decoder Message
message =
    Decode.map2 Message
        (field "posted" string)
        (field "content" string)


messages : Decoder (List Message)
messages =
    field "messages" (list message)


messagesQuery : Operation Query Named
messagesQuery =
    GraphQl.named "MessagesQuery"
        [ GraphQl.field "messages"
            |> GraphQl.withSelectors
                [ GraphQl.field "posted"
                , GraphQl.field "content"
                ]
        ]


messagesRequest : Operation Query Named -> Decoder (List Message) -> GraphQl.Request Query Named (List Message)
messagesRequest =
    GraphQl.query "/graphql"


query : (Result Http.Error (List Message) -> msg) -> Cmd msg
query handler =
    GraphQl.send handler (messagesRequest messagesQuery messages)
