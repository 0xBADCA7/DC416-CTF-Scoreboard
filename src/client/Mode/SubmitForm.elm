module Mode.SubmitForm exposing (Input(..), view, mutation)

-- Core packages

import Http
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput, onClick)
import Json.Decode as Decode exposing (Decoder, field, int, string, list)


-- Third-party packages

import GraphQl exposing (Operation, Mutation, Named)


-- Local packages

import Mode.Scoreboard exposing (Scoreboard, scoreboard)


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


submitMutation : String -> String -> Operation Mutation Named
submitMutation token flag =
    GraphQl.named "SubmitFlagMutation"
        [ GraphQl.field "submitFlag"
            |> GraphQl.withArgument "submissionToken" (GraphQl.string token)
            |> GraphQl.withArgument "flag" (GraphQl.string flag)
            |> GraphQl.withSelectors
                [ GraphQl.field "rank"
                , GraphQl.field "name"
                , GraphQl.field "score"
                , GraphQl.field "lastSubmission"
                ]
        ]


submitRequest : Operation Mutation Named -> Decoder Scoreboard -> GraphQl.Request Mutation Named Scoreboard
submitRequest =
    GraphQl.mutation "/graphql"


mutation : String -> String -> (Result Http.Error Scoreboard -> msg) -> Cmd msg
mutation token flag handler =
    GraphQl.send handler (submitRequest (submitMutation token flag) scoreboard)
