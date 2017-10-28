module Mode.Scoreboard exposing (Scoreboard(..), Team, view, query, scoreboard, team)

-- Core packages

import Http
import Html exposing (..)
import Html.Attributes exposing (id)
import Json.Decode as Decode exposing (Decoder, field, int, string, list)


-- Third-party packages

import GraphQl exposing (Operation, Query, Named)


-- MODEL


type Scoreboard
    = Scoreboard (List Team)


type alias Team =
    { rank : Int
    , name : String
    , score : Int
    , lastSubmission : String
    }



-- VIEW


view : Scoreboard -> Html msg
view (Scoreboard teams) =
    table [ id "scoreboard" ] <|
        List.append
            [ tr [ id "scoreboardHeader" ]
                [ th [] [ text "Rank" ]
                , th [] [ text "Name" ]
                , th [] [ text "Score" ]
                , th [] [ text "Last Submission" ]
                ]
            ]
            (List.map viewTeam teams)


viewTeam : Team -> Html msg
viewTeam team =
    tr []
        [ td [] [ text (toString team.rank) ]
        , td [] [ text team.name ]
        , td [] [ text (toString team.score) ]
        , td [] [ text team.lastSubmission ]
        ]



-- QUERY


team : Decoder Team
team =
    Decode.map4 Team
        (field "rank" int)
        (field "name" string)
        (field "score" int)
        (field "lastSubmission" string)


scoreboard : Decoder Scoreboard
scoreboard =
    Decode.map Scoreboard (field "teams" (list team))


scoreboardQuery : Operation Query Named
scoreboardQuery =
    GraphQl.named "ScoreboardQuery"
        [ GraphQl.field "teams"
            |> GraphQl.withSelectors
                [ GraphQl.field "rank"
                , GraphQl.field "name"
                , GraphQl.field "score"
                , GraphQl.field "lastSubmission"
                ]
        ]


scoreboardRequest : Operation Query Named -> Decoder Scoreboard -> GraphQl.Request Query Named Scoreboard
scoreboardRequest =
    GraphQl.query "/graphql"


query : (Result Http.Error Scoreboard -> msg) -> Cmd msg
query handler =
    GraphQl.send handler (scoreboardRequest scoreboardQuery scoreboard)
