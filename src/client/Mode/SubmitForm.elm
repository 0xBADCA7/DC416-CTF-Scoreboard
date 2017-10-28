module Mode.SubmitForm exposing (Input(..), SubmitFlagResponse, SubmitResponse, view, mutation)

-- Core packages

import Http
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput, onClick)
import Json.Decode as Decode exposing (Decoder, field, int, string, list, bool, maybe)


-- Third-party packages

import GraphQl exposing (Operation, Mutation, Named)


-- Local packages

import Mode.Scoreboard exposing (Scoreboard, Team, team, scoreboard)


-- VIEW


type Input
    = SubmissionToken String
    | Flag String
    | Submit


view : (Input -> msg) -> Html msg
view msg =
    div [ class "row" ]
        [ Html.form [ id "submitFlagForm", class "col s12" ]
            [ div [ class "row" ]
                [ div [ class "input-field col s12" ]
                    [ input
                        [ id "submissionTokenField"
                        , type_ "text"
                        , placeholder "submission token"
                        , onInput (\s -> (msg (SubmissionToken s)))
                        ]
                        []
                    ]
                ]
            , div [ class "row" ]
                [ div [ class "input-field col s12" ]
                    [ input
                        [ id "submissionFlagField"
                        , type_ "text"
                        , placeholder "flag"
                        , onInput (\s -> (msg (Flag s)))
                        ]
                        []
                    ]
                ]
            , a
                [ id "submitFlagBtn"
                , class "btn btnPrimary"
                , onClick (msg Submit)
                ]
                [ text "submit" ]
            ]
        ]



-- QUERY


type alias SubmitResponse =
    { correct : Bool
    , teams : List Team
    }


type alias SubmitFlagResponse =
    { submitFlag : SubmitResponse
    }


submitResponse : Decoder SubmitResponse
submitResponse =
    Decode.map2 SubmitResponse
        (field "correct" bool)
        (field "teams" (Decode.list team))


submitFlag : Decoder SubmitFlagResponse
submitFlag =
    Decode.map SubmitFlagResponse (field "submitFlag" submitResponse)


submitMutation : String -> String -> Operation Mutation Named
submitMutation token flag =
    GraphQl.named "SubmitFlagMutation"
        [ GraphQl.field "submitFlag"
            |> GraphQl.withArgument "submissionToken" (GraphQl.string token)
            |> GraphQl.withArgument "flag" (GraphQl.string flag)
            |> GraphQl.withSelectors
                [ GraphQl.field "correct"
                , GraphQl.field "teams"
                    |> GraphQl.withSelectors
                        [ GraphQl.field "rank"
                        , GraphQl.field "name"
                        , GraphQl.field "score"
                        , GraphQl.field "lastSubmission"
                        ]
                ]
        ]


submitRequest : Operation Mutation Named -> Decoder SubmitFlagResponse -> GraphQl.Request Mutation Named SubmitFlagResponse
submitRequest =
    GraphQl.mutation "/graphql"


mutation : String -> String -> (Result Http.Error SubmitFlagResponse -> msg) -> Cmd msg
mutation token flag handler =
    GraphQl.send handler (submitRequest (submitMutation token flag) submitFlag)
