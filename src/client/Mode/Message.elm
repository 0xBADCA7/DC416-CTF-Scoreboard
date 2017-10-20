module Mode.Message exposing (Message, view)

import Html exposing (..)
import Html.Attributes exposing (..)


type alias Message =
    { posted : String
    , content : String
    }


view : List Message -> Html msg
view messages =
    ul [ id "messageList", style [ ( "list-style", "none" ) ] ] <| List.map viewMessage <| List.reverse messages


viewMessage : Message -> Html msg
viewMessage { posted, content } =
    li [ class "adminMessage" ]
        [ span [ class "gray-text text-darken-4" ] [ text ("Posted " ++ posted) ]
        , p [] [ text content ]
        ]
