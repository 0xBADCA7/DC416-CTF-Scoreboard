module Mode.Scoreboard exposing (Scoreboard(..), Team, view, query, scoreboard, team)

-- Core packages

import Http
import Html exposing (..)
import Html.Attributes exposing (id)
import Json.Decode as Decode exposing (Decoder, maybe, field, int, string, list)


-- Third-party packages

import GraphQl exposing (Operation, Query, Named)


-- Local packages

import Util


-- MODEL


type Scoreboard
    = Scoreboard (List Team)


type alias Team =
    { name : String
    , score : Int
    , lastSubmission : Maybe Int
    }


compareTeams : Team -> Team -> Order
compareTeams team1 team2 =
    let
        maybeInt value =
            case value of
                Just n ->
                    n

                Nothing ->
                    0

        scoreOrder =
            compare team2.score team1.score

        submitOrder =
            compare (maybeInt team1.lastSubmission) (maybeInt team2.lastSubmission)
    in
        case ( scoreOrder, submitOrder ) of
            ( EQ, cmp ) ->
                cmp

            ( cmp, _ ) ->
                cmp



-- VIEW


view : Scoreboard -> Html msg
view (Scoreboard teams) =
    let
        sortedTeams =
            List.sortWith compareTeams teams

        ranks =
            List.range 1 (List.length sortedTeams)

        _ =
            Debug.log "Sorted teams" sortedTeams
    in
        table [ id "scoreboard" ] <|
            List.append
                [ tr [ id "scoreboardHeader" ]
                    [ th [] [ text "Rank" ]
                    , th [] [ text "Name" ]
                    , th [] [ text "Score" ]
                    , th [] [ text "Last Submission" ]
                    ]
                ]
                (List.map2 viewTeam ranks sortedTeams)


viewTeam : Int -> Team -> Html msg
viewTeam rank team =
    let
        submissionStr =
            case team.lastSubmission of
                Just time ->
                    Util.prettyDate time

                Nothing ->
                    "No submissions yet."
    in
        tr []
            [ td [] [ text (toString rank) ]
            , td [] [ text team.name ]
            , td [] [ text (toString team.score) ]
            , td [] [ text submissionStr ]
            ]



-- QUERY


team : Decoder Team
team =
    Decode.map3 Team
        (field "name" string)
        (field "score" int)
        (maybe (field "lastSubmission" int))


scoreboard : Decoder Scoreboard
scoreboard =
    Decode.map Scoreboard (field "teams" (list team))


scoreboardQuery : Operation Query Named
scoreboardQuery =
    GraphQl.named "ScoreboardQuery"
        [ GraphQl.field "teams"
            |> GraphQl.withSelectors
                [ GraphQl.field "name"
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
