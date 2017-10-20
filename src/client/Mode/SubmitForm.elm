module Mode.SubmitForm exposing (view)

import Html exposing (..)
import Html.Attributes exposing (..)


view : Html msg
view =
    div [ class "row" ]
        [ Html.form [ id "submitFlagForm", class "col s12" ]
            [ div [ class "row" ]
                [ div [ class "input-field col s12" ]
                    [ input
                        [ id "submissionTokenField"
                        , type_ "text"
                        , placeholder "submission token"
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
                        ]
                        []
                    ]
                ]
            , a [ id "submitFlagBtn", class "btn btnPrimary" ] [ text "submit" ]
            ]
        ]
