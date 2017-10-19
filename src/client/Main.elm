module Main exposing (main)

-- Standard imports

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)


-- Local imports

import Mode.Scoreboard as Scoreboard exposing (Scoreboard(..))


-- MAIN


main : Platform.Program Basics.Never Model Msg
main =
    Html.program
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        }



-- MODEL


type ViewMode
    = ScoreboardView
    | SubmitForm
    | MessagesView


type alias Model =
    { scoreboard : Scoreboard
    , submissionToken : String
    , submissionFlag : String
    , messages : List String
    , mode : ViewMode
    }


init : ( Model, Cmd Msg )
init =
    ( Model testScoreboard "" "" [] ScoreboardView, Cmd.none )


testScoreboard : Scoreboard
testScoreboard =
    Scoreboard
        [ { rank = 1, name = "Team one", score = 150, lastSubmission = "9:30" }
        , { rank = 2, name = "H4xx0R", score = 132, lastSubmission = "9:25" }
        , { rank = 3, name = "31337", score = 130, lastSubmission = "9:46" }
        , { rank = 4, name = "b4d455", score = 80, lastSubmission = "9:15" }
        , { rank = 5, name = "CTF TO", score = 68, lastSubmission = "9:08" }
        , { rank = 6, name = "T.", score = 35, lastSubmission = "8:12" }
        , { rank = 7, name = "DC416", score = 10, lastSubmission = "7:50" }
        ]



-- UPDATE


type Msg
    = SwitchMode ViewMode


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        SwitchMode mode ->
            ( { model | mode = mode }, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Html Msg
view model =
    div []
        [ viewNav model
        , viewMode model
        ]


viewNav : Model -> Html Msg
viewNav model =
    let
        navLinks =
            case model.mode of
                ScoreboardView ->
                    [ li [ class "active" ] [ a [ href "#", onClick (SwitchMode ScoreboardView) ] [ text "Scoreboard" ] ]
                    , li [] [ a [ href "#", onClick (SwitchMode MessagesView) ] [ text "Messages" ] ]
                    , li [] [ a [ href "#", onClick (SwitchMode SubmitForm) ] [ text "Submit" ] ]
                    ]

                SubmitForm ->
                    [ li [] [ a [ href "#", onClick (SwitchMode ScoreboardView) ] [ text "Scoreboard" ] ]
                    , li [ class "active" ] [ a [ href "#", onClick (SwitchMode MessagesView) ] [ text "Messages" ] ]
                    , li [] [ a [ href "#", onClick (SwitchMode SubmitForm) ] [ text "Submit" ] ]
                    ]

                MessagesView ->
                    [ li [] [ a [ href "#", onClick (SwitchMode ScoreboardView) ] [ text "Scoreboard" ] ]
                    , li [] [ a [ href "#", onClick (SwitchMode MessagesView) ] [ text "Messages" ] ]
                    , li [ class "active" ] [ a [ href "#", onClick (SwitchMode SubmitForm) ] [ text "Submit" ] ]
                    ]
    in
        nav []
            [ div [ class "mainContent nav-wrapper" ]
                [ a [ href "#", class "brand-logo" ] [ text "Scoreboard" ]
                , ul [ id "nav-mobile", class "right" ] navLinks
                ]
            ]


viewMode : Model -> Html Msg
viewMode model =
    let
        viewContent =
            case model.mode of
                ScoreboardView ->
                    [ span [ class "card-title gray-text text-darken-4" ] [ text "Scoreboard" ]
                    , Scoreboard.view model.scoreboard
                    ]

                SubmitForm ->
                    [ span [ class "card-title gray-text text-darken-4" ] [ text "Submit a flag" ]
                    , viewSubmitForm <| SwitchMode ScoreboardView
                    ]

                _ ->
                    []
    in
        div [ class "mainContent" ]
            [ div [ id "content", class "card large waves-effect waves-light" ]
                [ div [ class "card-content" ] viewContent
                ]
            ]


viewSubmitForm : Msg -> Html Msg
viewSubmitForm msg =
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
