module Mode.Message exposing (Message, view, query)

-- Core packages

import Time
import Date
import Http
import Html exposing (..)
import Html.Attributes exposing (id, style, class)
import Json.Decode as Decode exposing (Decoder, field, string, list, int)


-- Third-party packages

import GraphQl exposing (Operation, Query, Named)


-- MODEL


type alias Message =
    { posted : Int
    , content : String
    }



-- VIEW


view : List Message -> Html msg
view messages =
    let
        msgs =
            messages
                |> List.reverse
                |> List.map viewMessage
    in
        ul
            [ id "messageList"
            , style [ ( "list-style", "none" ) ]
            ]
            msgs


viewMessage : Message -> Html msg
viewMessage { posted, content } =
    let
        postedDate =
            posted
                |> toFloat
                |> (\x -> x * Time.second)
                |> Date.fromTime

        dateStr =
            (toString <| Date.month postedDate) ++ " " ++ (toString <| Date.day postedDate)

        hourStr =
            postedDate
                |> Date.hour
                |> toString
                |> String.padLeft 2 '0'

        minuteStr =
            postedDate
                |> Date.minute
                |> toString
                |> String.padLeft 2 '0'

        postedStr =
            "Posted on " ++ dateStr ++ " at " ++ hourStr ++ ":" ++ minuteStr
    in
        li [ class "adminMessage" ]
            [ span [ class "gray-text text-darken-4" ] [ text postedStr ]
            , p [] [ text content ]
            ]



-- QUERY


message : Decoder Message
message =
    Decode.map2 Message
        (field "posted" int)
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
