module Action exposing (Action(..), Team, Scoreboard, fetchScoreboard, viewScoreboard)

-- Core imports

import Html exposing (..)
import Html.Attributes exposing (id)
import Json.Decode exposing (Decoder, int, string, list)


-- Package imports

import Http
import Json.Decode.Pipeline exposing (decode, required)


{-| Represents different states the client model can be moved into, from which we would like to run commands.
-}
type Action
    = Nothing
    | ViewScoreboard


type alias Team =
    { name : String
    , score : Int
    , position : Int
    , members : String
    , lastSubmission : String
    }


type alias Scoreboard =
    { teams : List Team
    }



-- Functions for executing the behaviour we want actions to result in.


fetchScoreboard : Http.Request Scoreboard
fetchScoreboard =
    Http.get "/teams/scoreboard" scoreboard



-- JSON encoders and decoders.


team : Decoder Team
team =
    decode Team
        |> required "name" string
        |> required "score" int
        |> required "position" int
        |> required "members" string
        |> required "lastSubmission" string


scoreboard : Decoder Scoreboard
scoreboard =
    decode Scoreboard
        |> required "teams" (list team)



-- View functions.


sbHeaderRow : Html msg
sbHeaderRow =
    tr [ id "scoreheader" ]
        [ th [] [ text "Place" ]
        , th [] [ text "Team Name" ]
        , th [] [ text "Score" ]
        , th [] [ text "Members" ]
        , th [] [ text "Last Submission" ]
        ]


viewTeam : Team -> Html msg
viewTeam t =
    tr []
        [ td [] [ text <| toString t.position ]
        , td [] [ text t.name ]
        , td [] [ text <| toString t.score ]
        , td [] [ text <| t.members ]
        , td [] [ text <| t.lastSubmission ]
        ]


viewScoreboard : Scoreboard -> Html msg
viewScoreboard sb =
    case sb.teams of
        [] ->
            table [ id "scoreboard" ] [ sbHeaderRow ]

        teams ->
            let
                sortedTeams =
                    List.sortBy .position teams
            in
                table [ id "scoreboard" ] <| sbHeaderRow :: List.map viewTeam sortedTeams
